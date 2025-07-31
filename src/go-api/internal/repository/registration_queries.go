package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func FindRegistrationByToken(ctx context.Context, db sqlx.QueryerContext, token string) (*model.ExtendedRegistration, error) {
	regs, err := listRegistrations(ctx, db, "r.token = $1 IS NULL", token)
	if err != nil {
		return nil, err
	}
	if len(regs) == 0 {
		return nil, errors.WithStack(sql.ErrNoRows)
	}
	return &regs[0], nil
}

func ListRegistrations(ctx context.Context, db sqlx.QueryerContext) ([]model.ExtendedRegistration, error) {
	return listRegistrations(ctx, db, "r.deleted_at IS NULL AND s.deleted_at IS NULL")
}

func ListEventRegistrations(ctx context.Context, db sqlx.QueryerContext, eventID int) ([]model.ExtendedRegistration, error) {
	return listRegistrations(ctx, db, "e.event_id=$1 AND deleted_at IS NULL", eventID)
}

func ListRegistrationsWithoutNotification(ctx context.Context, db sqlx.QueryerContext) ([]model.ExtendedRegistration, error) {
	return listRegistrations(ctx, db, "r.notification_sent_at IS NULL AND r.deleted_at IS NULL AND s.deleted_at IS NULL")
}

func FindRegistrationByID(ctx context.Context, db sqlx.QueryerContext, regID int) (*model.ExtendedRegistration, error) {
	regs, err := listRegistrations(ctx, db, "r.id = $1 AND r.deleted_at IS NULL", regID)
	if err != nil {
		return nil, err
	}
	if len(regs) == 0 {
		return nil, errors.WithStack(sql.ErrNoRows)
	}
	return &regs[0], nil
}

func UpdateRegistration(ctx context.Context, tx *sqlx.Tx, reg *model.Registration) error {
	stmt, err := tx.PrepareNamedContext(ctx, `
		UPDATE registrations SET 
		    amount = :amount,
			payed = :payed,
			admin_note = :admin_note,
			updated_at = NOW()
		WHERE id = :id
		RETURNING id
	`)
	if err != nil {
		return errors.Wrap(err, "failed to prepare query")
	}

	var updated int
	if err := stmt.GetContext(ctx, &updated, reg); err != nil {
		return errors.Wrap(err, "failed to update a registration")
	}
	return nil
}

func MarkRegistrationAsNotified(ctx context.Context, tx sqlx.QueryerContext, id int) error {
	var notifiedID int
	err := sqlx.GetContext(ctx, tx, &notifiedID, `
		UPDATE registrations SET 
			notification_sent_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id
	`, id)
	if err != nil {
		return errors.Wrap(err, "failed to mark registration as notified")
	}
	return nil
}

func DeleteRegistrationByID(ctx context.Context, db sqlx.QueryerContext, regID int) (*model.Registration, error) {
	return markRegistrationAsDeleted(ctx, db, regID)
}

func markRegistrationAsDeleted(ctx context.Context, db sqlx.QueryerContext, id int) (*model.Registration, error) {
	var deleted model.Registration
	err := sqlx.GetContext(ctx, db, &deleted, `
		UPDATE registrations SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		RETURNING *
	`, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to mark registration as deleted")
	}
	return &deleted, nil
}

func listRegistrations(ctx context.Context, db sqlx.QueryerContext, where string, args ...interface{}) ([]model.ExtendedRegistration, error) {
	const queryTemplate = `
		SELECT
			r.id,
			r.name,
			r.surname,
			e.id AS event_id,
			e.title,
			json_agg(json_build_object('id', d.id, 'description', d.description) ORDER BY d.id)  AS days,
			r.gender,
			r.date_of_birth,
			r.finished_school,
			r.attended_previous,
			r.city,
			r.pills,
			r.notes,
			r.parent_name,
			r.attended_activities,
			r.problems,
			r.parent_surname,
			r.email,
			r.phone,
			r.amount,
			r.payed,
			r.admin_note,
			r.created_at,
			r.token,
			r.updated_at,
			r.discount,
			r.promo_code,
			r.notification_sent_at,
			r.specific_symbol
		FROM registrations r
		LEFT JOIN signups s ON r.id = s.registration_id
		LEFT JOIN days d ON s.day_id = d.id
		LEFT JOIN events e ON d.event_id = e.id
		%s
		GROUP BY r.id, r.name, r.surname, e.title, e.id, r.gender, r.date_of_birth,
			r.finished_school, r.attended_previous, r.city, r.pills, r.notes,
			r.parent_name,  r.attended_activities, r.problems, r.parent_surname,
			r.email,  r.phone , r.amount, r.payed, r.created_at,
			r.updated_at`

	var condition string
	if where != "" {
		condition = fmt.Sprintf("WHERE %s", where)
	}

	var res []model.ExtendedRegistration
	if err := sqlx.SelectContext(ctx, db, &res, fmt.Sprintf(queryTemplate, condition), args...); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func CreateRegistration(ctx context.Context, db sqlx.QueryerContext, r model.Registration) (*model.Registration, error) {
	var reg model.Registration
	err := sqlx.GetContext(ctx, db, &reg, `
			INSERT INTO registrations(
				name,
				surname,
				gender,
				amount,
				token,
				date_of_birth,
				finished_school,
				attended_previous,
				city,
				pills,
				notes,
				parent_name,
				parent_surname,
				email,
				phone,
				attended_activities,
				problems,
				promo_code,
			    discount,
				created_at,
				updated_at
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12,
				$13,
				$14,
				$15,
				$16,
				$17,
				$18,
			    $19,
				NOW(),
				NOW()
			) RETURNING *
		`, r.Name, r.Surname, r.Gender, r.Amount, r.Token,
		r.DateOfBirth, r.FinishedSchool, r.AttendedPrevious, r.City,
		r.Pills, r.Notes, r.ParentName, r.ParentSurname, r.Email,
		r.Phone, r.AttendedActivities, r.Problems, r.PromoCode, r.Discount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a registration")
	}
	return &reg, nil
}

func CreateSignup(ctx context.Context, db sqlx.QueryerContext, s model.Signup) (*model.Signup, error) {
	var res model.Signup
	err := sqlx.GetContext(ctx, db, &res, `
			INSERT INTO signups(
				day_id,
				registration_id,
				state,
				created_at,
				updated_at
			) VALUES (
				$1,
				$2,
				$3,
				NOW(),
				NOW()
			) RETURNING *
		`, s.DayID, s.RegistrationID, s.State)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a signup entry")
	}
	return &res, nil
}
