package repository

import (
	"context"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func CreatePromoCode(ctx context.Context, db sqlx.QueryerContext, data model.PromoCode) (*model.PromoCode, error) {
	var promoCode model.PromoCode
	err := sqlx.GetContext(ctx, db, &promoCode, `
		INSERT INTO promo_codes(
		    email,
		    key,
		    available_registrations,
		    updated_at,
			created_at
		) VALUES(
				$1,
				$2,
				$3,
				NOW(),
				NOW() 
		) RETURNING *
	`, data.Email, data.Key, data.AvailableRegistrations)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create promo code")
	}
	return &promoCode, nil
}

func FindPromoCodeByKey(ctx context.Context, tx sqlx.QueryerContext, key string) (*model.PromoCode, error) {
	var code model.PromoCode
	if err := sqlx.GetContext(ctx, tx, &code, `SELECT * FROM promo_codes WHERE key = $1`, key); err != nil {
		return nil, errors.Wrap(err, "failed to find promo_code")
	}
	return &code, nil
}

func DecrementAvailableRegistrationsPromoCodeByKey(ctx context.Context, tx sqlx.QueryerContext, key string) (*model.PromoCode, error) {
	var code model.PromoCode
	if err := sqlx.GetContext(ctx, tx, &code, `UPDATE promo_codes SET available_registrations=available_registrations - 1 WHERE key = $1 RETURNING *`, key); err != nil {
		return nil, errors.Wrap(err, "failed to mark promo_code as used")
	}
	return &code, nil
}
