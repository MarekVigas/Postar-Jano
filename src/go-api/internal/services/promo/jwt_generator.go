package promo

import (
	"context"
	"github.com/MarekVigas/Postar-Jano/internal/repository"
	"github.com/MarekVigas/Postar-Jano/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/MarekVigas/Postar-Jano/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	audience = "Leto 2024"
	issuer   = "sbb.sk"
)

type JWTGenerator struct {
	logger         *zap.Logger
	promoSecret    []byte
	activationDate *time.Time
	expirationDate *time.Time
}

func NewJWTGenerator(logger *zap.Logger, secret []byte, activationDate *time.Time, expirationDate *time.Time) *JWTGenerator {
	logger.Debug("New promo generator created", zap.Timep("activation_date", activationDate),
		zap.Timep("expiration_date", expirationDate))
	return &JWTGenerator{
		logger:         logger,
		promoSecret:    secret,
		activationDate: activationDate,
		expirationDate: expirationDate,
	}
}

func (g *JWTGenerator) GenerateToken(ctx context.Context, db sqlx.QueryerContext, email string, registrationCount int) (token string, err error) {
	key := uuid.New().String()
	claims := jwt.RegisteredClaims{
		Audience: jwt.ClaimStrings{audience},
		ID:       key,
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Issuer:   issuer,
	}
	if g.activationDate != nil {
		claims.NotBefore = jwt.NewNumericDate(g.activationDate.UTC())
	}
	if g.expirationDate != nil {
		claims.ExpiresAt = jwt.NewNumericDate(g.expirationDate.UTC())
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(g.promoSecret)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if _, err := repository.CreatePromoCode(ctx, db, model.PromoCode{
		Email:                  email,
		Key:                    key,
		AvailableRegistrations: registrationCount,
	}); err != nil {
		return "", err
	}
	return token, nil
}

func (g *JWTGenerator) ValidateToken(ctx context.Context, tx sqlx.QueryerContext, token string) (code *model.PromoCode, err error) {
	var standardClaims jwt.RegisteredClaims
	decodedToken, err := jwt.ParseWithClaims(token, &standardClaims, func(token *jwt.Token) (interface{}, error) {
		return g.promoSecret, nil
	})
	if err != nil {
		return nil, errors.WithStack(ErrInvalid)
	}
	if !decodedToken.Valid {
		return nil, errors.WithStack(ErrInvalid)
	}

	promoCode, err := repository.FindPromoCodeByKey(ctx, tx, standardClaims.ID)
	if err != nil {
		logger.UnexpectedError(ctx, err).Error("Failed to find promo code", zap.String("token", token))
		return nil, err
	}
	if promoCode.AvailableRegistrations <= 0 {
		return nil, errors.WithStack(ErrAlreadyUsed)
	}
	return promoCode, nil
}
