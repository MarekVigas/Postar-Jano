package model

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PromoCode struct {
	ID                     int       `json:"id" db:"id"`
	Email                  string    `json:"email" db:"email"`
	Key                    string    `json:"key" db:"key"`
	AvailableRegistrations int       `json:"available_registrations" db:"available_registrations"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
}

func (m *PromoCode) Create(ctx context.Context, db sqlx.QueryerContext) (*PromoCode, error) {
	var promoCode PromoCode
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
	`, m.Email, m.Key, m.AvailableRegistrations)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create promo code")
	}
	return &promoCode, nil
}

func FindPromoCodeByKey(ctx context.Context, tx sqlx.QueryerContext, key string) (*PromoCode, error) {
	var code PromoCode
	if err := sqlx.GetContext(ctx, tx, &code, `SELECT * FROM promo_codes WHERE key = $1`, key); err != nil {
		return nil, errors.Wrap(err, "failed to find promo_code")
	}
	return &code, nil
}

func DecrementAvailableRegistrationsPromoCodeByKey(ctx context.Context, tx sqlx.QueryerContext, key string) (*PromoCode, error) {
	var code PromoCode
	if err := sqlx.GetContext(ctx, tx, &code, `UPDATE promo_codes SET available_registrations=available_registrations - 1 WHERE key = $1 RETURNING *`, key); err != nil {
		return nil, errors.Wrap(err, "failed to mark promo_code as used")
	}
	return &code, nil
}
