package model

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Day struct {
	ID          int    `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
	Capacity    int    `json:"capacity" db:"capacity"`
	LimitBoys   *int   `json:"limit_boys" db:"limit_boys"`
	LimitGirls  *int   `json:"limit_girls" db:"limit_girls"`
	Price       int    `json:"price" db:"price"`
	EventID     int    `json:"-" db:"event_id"`
}

func (d *Day) Create(ctx context.Context, db sqlx.ExtContext) error {
	return errors.WithStack(sqlx.GetContext(ctx, db, d, `
		INSERT INTO days(
			description,
			capacity,
			limit_boys,
			limit_girls,
			price,
			event_id
		) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`, d.Description, d.Capacity, d.LimitBoys, d.LimitGirls, d.Price, d.EventID))
}
