package repository

import (
	"context"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func ListEvents(ctx context.Context, tx *sqlx.Tx) ([]model.Event, error) {
	var events []model.Event
	if err := sqlx.SelectContext(ctx, tx, &events, `
			SELECT 
				ev.id,
				ev.title,
				ev.description,
				ev.date_from,
				ev.date_to,
				ev.location,
				ev.min_age,
				ev.max_age,
				ev.info,
				ev.photo,
				ev.time,
				ev.price,
				ev.mail_info,
				ev.active,
				ev.promo_registration,
				ev.iban,
				ev.payment_reference,
				ev.promo_discount,
				o.id AS owner_id,
				o.name AS owner_name,
				o.surname AS owner_surname,
				o.email AS owner_email,
				o.phone AS owner_phone,
				o.photo AS owner_photo,
				o.gender AS owner_gender
			FROM events ev
			LEFT JOIN owners o ON o.id = ev.owner_id
		`); err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range events {
		if err := sqlx.SelectContext(ctx, tx, &events[i].Days, `
			SELECT 
				id,
				capacity,
				limit_boys,
				limit_girls,
				description,
				price
			FROM days
			WHERE event_id = $1
			ORDER BY id
		`, events[i].ID); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return events, nil
}

func FindEvent(ctx context.Context, tx *sqlx.Tx, id int) (*model.Event, error) {
	var event model.Event
	if err := sqlx.GetContext(ctx, tx, &event, `
			SELECT 
				ev.id,
				ev.title,
				ev.description,
				ev.date_from,
				ev.date_to,
				ev.location,
				ev.min_age,
				ev.max_age,
				ev.info,
				ev.photo,
				ev.time,
				ev.price,
				ev.mail_info,
				ev.active,
				ev.promo_registration,
				ev.iban,
				ev.payment_reference,
				ev.promo_discount,
				o.id AS owner_id,
				o.name AS owner_name,
				o.surname AS owner_surname,
				o.email AS owner_email,
				o.phone AS owner_phone,
				o.photo AS owner_photo,
				o.gender AS owner_gender
			FROM events ev
			LEFT JOIN owners o ON o.id = ev.owner_id
			WHERE ev.id = $1
		`, id); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := sqlx.SelectContext(ctx, tx, &event.Days, `
			SELECT 
				id,
				capacity,
				limit_boys,
				limit_girls,
				description,
				price
			FROM days
			WHERE event_id = $1
			ORDER BY id
		`, id); err != nil {
		return nil, errors.WithStack(err)
	}

	return &event, nil
}

func CreateDay(ctx context.Context, db sqlx.QueryerContext, d model.Day) (*model.Day, error) {
	var result model.Day
	if err := sqlx.GetContext(ctx, db, &result, `
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
	`, d.Description, d.Capacity, d.LimitBoys, d.LimitGirls, d.Price, d.EventID); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func CreateEvent(ctx context.Context, db sqlx.ExtContext, e model.Event) (*model.Event, error) {
	var result model.Event
	if err := sqlx.GetContext(ctx, db, &result, `
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
			owner_id,
		    iban,
		    payment_reference,
		    promo_discount
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING *
	`, e.Title, e.Description, e.DateFrom, e.DateTo, e.Location, e.MinAge, e.MaxAge,
		e.Info, e.Photo, e.Time, e.Price, e.MailInfo, e.Active, e.OwnerID, e.IBAN,
		e.PaymentReference, e.PromoDiscount); err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func CreateOwner(ctx context.Context, db sqlx.ExtContext, o model.Owner, plainTextPass string) (*model.Owner, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(plainTextPass), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var owner model.Owner
	if err := sqlx.GetContext(ctx, db, &owner, `
		INSERT INTO owners (
			name,
			surname,
			gender,
			username,
			pass,
			email,
			phone,
			photo
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *
	`, o.Name, o.Surname, o.Gender, o.Username, pass, o.Email,
		o.Phone, o.Photo); err != nil {
		return nil, errors.WithStack(err)
	}
	return &owner, nil
}
