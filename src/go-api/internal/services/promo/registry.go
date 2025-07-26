package promo

import (
	"context"
	"database/sql"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/repository"
	"github.com/MarekVigas/Postar-Jano/internal/resources"
	"github.com/MarekVigas/Postar-Jano/internal/services/mailer/templates"
	"github.com/MarekVigas/Postar-Jano/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Generator interface {
	GenerateToken(ctx context.Context, db sqlx.QueryerContext, email string, registrationCount int) (token string, err error)
	ValidateToken(ctx context.Context, db sqlx.QueryerContext, token string) (*model.PromoCode, error)
}

type EmailSender interface {
	PromoMail(ctx context.Context, req *templates.PromoReq) error
}

type Registry struct {
	postgresDB  *repository.PostgresDB
	generator   Generator
	emailSender EmailSender
}

func NewRegistry(postgresDB *repository.PostgresDB, generator Generator, sender EmailSender) *Registry {
	return &Registry{
		postgresDB:  postgresDB,
		generator:   generator,
		emailSender: sender,
	}
}

func (r *Registry) GenerateToken(ctx context.Context, email string, registrationCount int, sendEmail bool) (string, error) {
	// Generate token
	token, err := r.generator.GenerateToken(ctx, r.postgresDB.QueryerContext(), email, registrationCount)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to generate promo code.", zap.Error(err))
		return "", err
	}
	if sendEmail {
		if err := r.emailSender.PromoMail(ctx, &templates.PromoReq{
			Mail:                   email,
			Token:                  token,
			AvailableRegistrations: registrationCount,
		}); err != nil {
			logger.FromCtx(ctx).Error("Failed to send a confirmation mail.", zap.Error(err))
			return "", err
		}
	}
	return token, nil
}

func (r *Registry) ValidateToken(ctx context.Context, token string) (*resources.ValidatePromoCodeResp, error) {
	promoCode, err := r.ValidateTokenWithQueryerContext(ctx, r.postgresDB.QueryerContext(), token)
	if err != nil {
		switch err := errors.Cause(err); {
		case errors.Is(err, sql.ErrNoRows), errors.Is(err, ErrAlreadyUsed), errors.Is(err, ErrInvalid):
			return resources.InvalidPromoCodeResp(), nil
		default:
			logger.FromCtx(ctx).Error("Error during token validation.", zap.Error(err))
			return nil, err
		}
	}
	return resources.ValidPromoCodeResp(promoCode.AvailableRegistrations), nil
}

func (r *Registry) ValidateTokenWithQueryerContext(ctx context.Context, db sqlx.QueryerContext, token string) (*model.PromoCode, error) {
	promoCode, err := r.generator.ValidateToken(ctx, db, token)
	if err != nil {
		return nil, err
	}
	return promoCode, nil
}

func (r *Registry) MarkTokenUsage(ctx context.Context, tx sqlx.QueryerContext, key string) (err error) {
	if _, err := repository.DecrementAvailableRegistrationsPromoCodeByKey(ctx, tx, key); err != nil {
		logger.FromCtx(ctx).Error("Failed to decrement promo code by key", zap.String("key", key))
		return errors.WithStack(err)
	}
	return nil
}
