package promo

import (
	"context"
	"fmt"
	"regexp"

	"github.com/MarekVigas/Postar-Jano/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type SimpleGenerator struct {
	logger *zap.Logger
}

func NewSimpleGenerator(logger *zap.Logger) *JWTGenerator {
	logger.Debug("New simple promo generator created")
	return &JWTGenerator{
		logger: logger,
	}
}

func (g *SimpleGenerator) GenerateToken(ctx context.Context, tx *sqlx.Tx, email string, registrationCount int) (token string, err error) {
	key := uuid.New().String()
	if _, err := (&model.PromoCode{
		Email:                  email,
		Key:                    key,
		AvailableRegistrations: registrationCount,
	}).Create(ctx, tx); err != nil {
		return "", err
	}
	return key, nil
}

func (g *SimpleGenerator) ValidateToken(ctx context.Context, tx *sqlx.Tx, token string) (code *model.PromoCode, err error) {
	match, err := regexp.MatchString("^[a-zA-Z0-9_-]*", token)
	if err == nil {
		fmt.Println("Match:", match)
	}

	promoCode, err := model.FindPromoCodeByKey(ctx, tx, token)
	if err != nil {
		return nil, err
	}
	if promoCode.AvailableRegistrations <= 0 {
		return nil, errors.WithStack(ErrAlreadyUsed)
	}
	return promoCode, nil
}

func (g *SimpleGenerator) MarkTokenUsage(ctx context.Context, tx *sqlx.Tx, key string) (err error) {
	if _, err := model.DecrementAvailableRegistrationsPromoCodeByKey(ctx, tx, key); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
