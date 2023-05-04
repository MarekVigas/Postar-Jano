package model

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Event struct {
	ID                int     `json:"id" db:"id"`
	Title             string  `json:"title" db:"title"`
	Description       string  `json:"description" db:"description"`
	DateFrom          string  `json:"date_from" db:"date_from"`
	DateTo            string  `json:"date_to" db:"date_to"`
	Location          string  `json:"location" db:"location"`
	MinAge            int     `json:"min_age" db:"min_age"`
	MaxAge            int     `json:"max_age" db:"max_age"`
	Info              *string `json:"info" db:"info"`
	Photo             string  `json:"photo" db:"photo"`
	Time              *string `json:"time" db:"time"`
	Price             *string `json:"price" db:"price"`
	MailInfo          *string `json:"mail_info" db:"mail_info"`
	Active            bool    `json:"active" db:"active"`
	PromoRegistration bool    `json:"promo_registration" db:"promo_registration"`
	EventOwner        `json:"owner"`
	Days              []Day `json:"days"`
}

type EventOwner struct {
	OwnerID      int    `json:"id" db:"owner_id"`
	OwnerName    string `json:"name" db:"owner_name"`
	OwnerSurname string `json:"surname" db:"owner_surname"`
	OwnerEmail   string `json:"email" db:"owner_email"`
	OwnerPhone   string `json:"phone" db:"owner_phone"`
	OwnerPhoto   string `json:"photo" db:"owner_photo"`
	OwnerGender  string `json:"gender" db:"owner_gender"`
}

func (e *Event) Create(ctx context.Context, db sqlx.ExtContext) error {
	return errors.WithStack(sqlx.GetContext(ctx, db, e, `
		INSERT INTO events(
			title,
			description,
			date_from,
			date_to,
			location,
			min_age,
			max_age,
			info,
			photo,
			time,
			price,
			mail_info,
			active,
			owner_id
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING *
	`, e.Title, e.Description, e.DateFrom, e.DateTo, e.Location, e.MinAge, e.MaxAge,
		e.Info, e.Photo, e.Time, e.Price, e.MailInfo, e.Active, e.OwnerID))
}
