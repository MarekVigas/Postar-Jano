package promo

import (
	"context"
	"time"

	"github.com/MarekVigas/Postar-Jano/internal/model"

	"github.com/dgrijalva/jwt-go"
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

func (g *JWTGenerator) GenerateToken(ctx context.Context, tx sqlx.QueryerContext, email string, registrationCount int) (token string, err error) {
	key := uuid.New().String()
	claims := jwt.StandardClaims{
		Audience: audience,
		Id:       key,
		IssuedAt: time.Now().UTC().Unix(),
		Issuer:   issuer,
	}
	if g.activationDate != nil {
		claims.NotBefore = g.activationDate.UTC().Unix()
	}
	if g.expirationDate != nil {
		claims.ExpiresAt = g.expirationDate.UTC().Unix()
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(g.promoSecret)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if _, err := (&model.PromoCode{
		Email:                  email,
		Key:                    key,
		AvailableRegistrations: registrationCount,
	}).Create(ctx, tx); err != nil {
		return "", err
	}
	return token, nil
}

func (g *JWTGenerator) ValidateToken(ctx context.Context, tx sqlx.QueryerContext, token string) (code *model.PromoCode, err error) {
	var standardClaims jwt.StandardClaims
	decodedToken, err := jwt.ParseWithClaims(token, &standardClaims, func(token *jwt.Token) (interface{}, error) {
		return g.promoSecret, nil
	})
	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			return nil, errors.WithStack(ErrInvalid)
		}
		return nil, errors.WithStack(err)
	}
	if !decodedToken.Valid {
		return nil, errors.WithStack(ErrInvalid)
	}

	promoCode, err := model.FindPromoCodeByKey(ctx, tx, standardClaims.Id)
	if err != nil {
		return nil, err
	}
	if promoCode.AvailableRegistrations <= 0 {
		return nil, errors.WithStack(ErrAlreadyUsed)
	}
	return promoCode, nil
}

func (g *JWTGenerator) MarkTokenUsage(ctx context.Context, tx sqlx.QueryerContext, key string) (err error) {
	if _, err := model.DecrementAvailableRegistrationsPromoCodeByKey(ctx, tx, key); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
