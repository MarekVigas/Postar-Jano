package repository

import (
	"context"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
