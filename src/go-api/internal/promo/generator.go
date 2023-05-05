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
	audience = "Leto 2023"
	issuer   = "sbb.sk"
)

var (
	ErrAlreadyUsed = errors.New("already used")
)

type Generator struct {
	logger         *zap.Logger
	promoSecret    []byte
	activationDate *time.Time
	expirationDate *time.Time
}

func NewGenerator(logger *zap.Logger, secret []byte, activationDate *time.Time, expirationDate *time.Time) *Generator {
	logger.Debug("New promo generator created", zap.Timep("activation_date", activationDate),
		zap.Timep("expiration_date", expirationDate))
	return &Generator{
		logger:         logger,
		promoSecret:    secret,
		activationDate: activationDate,
		expirationDate: expirationDate,
	}
}

func (g *Generator) GenerateToken(ctx context.Context, tx *sqlx.Tx, email string, registrationCount int) (token string, err error) {
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

func (g *Generator) ValidateToken(ctx context.Context, tx *sqlx.Tx, token string) (code *model.PromoCode, err error) {
	var standardClaims jwt.StandardClaims
	decodedToken, err := jwt.ParseWithClaims(token, &standardClaims, func(token *jwt.Token) (interface{}, error) {
		return g.promoSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !decodedToken.Valid {
		return nil, errors.New("Invalid token!")
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

func (g *Generator) MarkTokenUsage(ctx context.Context, tx *sqlx.Tx, key string) (err error) {
	if _, err := model.DecrementAvailableRegistrationsPromoCodeByKey(ctx, tx, key); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
