package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/resources"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PostgresRepo struct {
	db           *sqlx.DB
	promoManager PromoManager
}

var (
	ErrOverLimit = errors.New("Limit exceeded")
	ErrNotActive = errors.New("Event not active")
)

type PromoManager interface {
	GenerateToken(ctx context.Context, tx sqlx.QueryerContext, email string, registrationCount int) (token string, err error)
	ValidateToken(ctx context.Context, tx sqlx.QueryerContext, token string) (code *model.PromoCode, err error)
	MarkTokenUsage(ctx context.Context, tx sqlx.QueryerContext, key string) (err error)
}

func NewPostgresRepo(db *sql.DB, manager PromoManager) *PostgresRepo {
	return &PostgresRepo{
		db:           sqlx.NewDb(db, "postgres"),
		promoManager: manager,
	}
}

func getAllTables() []string {
	return []string{
		"registrations",
		"signups",
		"owners",
		"events",
		"days",
		"promo_codes",
	}
}

func Reset(ctx context.Context, db sqlx.ExecerContext) error {
	for _, tableName := range getAllTables() {
		if _, err := db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (repo *PostgresRepo) Ping(ctx context.Context) error {
	if err := repo.db.PingContext(ctx); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (repo *PostgresRepo) ListEvents(ctx context.Context) ([]model.Event, error) {
	var events []model.Event
	err := repo.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
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
			return errors.WithStack(err)
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
				return errors.WithStack(err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (repo *PostgresRepo) FindEvent(ctx context.Context, id int) (*model.Event, error) {
	var event model.Event
	err := repo.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
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
			return errors.WithStack(err)
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
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (repo *PostgresRepo) Register(ctx context.Context, req *resources.RegisterReq, eventID int) (*model.RegResult, bool, error) {
	token, err := uuid.NewUUID()
	if err != nil {
		return nil, false, errors.WithStack(err)
	}

	event, err := repo.FindEvent(ctx, eventID)
	if err != nil {
		return nil, false, err
	}

	if !event.Active && !(event.PromoRegistration && req.PromoCode != nil) {
		return nil, false, ErrNotActive
	}

	res := model.RegResult{
		Token: token.String(),
		Event: event,
	}

	if err := repo.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var promoKey string
		if req.PromoCode != nil {
			// Validate token
			promoCode, err := repo.promoManager.ValidateToken(ctx, tx, *req.PromoCode)
			if err != nil {
				return err
			}
			promoKey = promoCode.Key
			// Decrement remaining registrations
			if err := repo.promoManager.MarkTokenUsage(ctx, tx, promoCode.Key); err != nil {
				return err
			}
		}

		// List stats.
		stats, err := repo.getStatInternal(ctx, tx, eventID)
		if err != nil {
			return err
		}

		// Validate limits
		for i := range req.DayIDs {
			for _, stat := range stats {
				if stat.DayID != req.DayIDs[i] {
					continue
				}
				if (stat.GirlsCount + stat.BoysCount) >= stat.Capacity {
					break
				}
				if req.Child.Gender == "male" && stat.LimitBoys != nil {
					if stat.BoysCount >= *stat.LimitBoys {
						break
					}
				}

				if req.Child.Gender == "female" && stat.LimitGirls != nil {
					if stat.GirlsCount >= *stat.LimitGirls {
						break
					}
				}
				res.RegisteredIDs = append(res.RegisteredIDs, stat.DayID)
			}
		}

		if len(req.DayIDs) != len(res.RegisteredIDs) {
			res.Success = false
			if len(res.RegisteredIDs) == 0 {
				res.RegisteredIDs = []int{}
			}
			return ErrOverLimit
		}
		res.Success = true

		// Compute price.
		var amount int
		for _, dID := range res.RegisteredIDs {
			for _, d := range event.Days {
				if dID != d.ID {
					continue
				}
				amount += d.Price
				res.RegisteredDesc = append(res.RegisteredDesc, d.Description)
				break
			}
		}

		// Apply promo discount
		var discount int
		if promoKey != "" {
			discount = event.PromoDiscount
		}

		newNullString := func(s string) *string {
			if s != "" {
				return &s
			}
			return nil
		}

		newNullInt := func(val int) *int {
			if val != 0 {
				return &val
			}
			return nil
		}

		// Insert into registrations.
		reg, err := (&model.Registration{
			Name:               req.Child.Name,
			Surname:            req.Child.Surname,
			Gender:             req.Child.Gender,
			DateOfBirth:        req.Child.DateOfBirth,
			FinishedSchool:     req.Child.FinishedSchool,
			AttendedPrevious:   req.Child.AttendedPrevious,
			AttendedActivities: req.Membership.AttendedActivities,
			City:               req.Child.City,
			Pills:              req.Medicine.Pills,
			Problems:           req.Health.Problems,
			Notes:              req.Notes,
			ParentName:         req.Parent.Name,
			ParentSurname:      req.Parent.Surname,
			Email:              req.Parent.Email,
			Phone:              req.Parent.Phone,
			Amount:             amount,
			Token:              token.String(),
			PromoCode:          newNullString(promoKey),
			Discount:           newNullInt(discount),
		}).Create(ctx, tx)
		if err != nil {
			return err
		}
		res.Reg = *reg

		// Insert into signups.
		for _, dayID := range req.DayIDs {
			_, err := tx.ExecContext(ctx, `
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
		`, dayID, res.Reg.ID, "init")
			if err != nil {
				return errors.Wrap(err, "failed to create a signup")
			}
		}

		return nil
	}); err != nil {
		if errors.Cause(err) == ErrOverLimit {
			return &res, false, nil
		}
		return nil, false, err
	}

	return &res, true, nil
}

func (repo *PostgresRepo) getStats(ctx context.Context, db sqlx.QueryerContext, where string, args ...interface{}) ([]model.Stat, error) {
	var stats []model.Stat
	if err := sqlx.SelectContext(ctx, db, &stats, fmt.Sprintf(`
			WITH
			boys AS (
				SELECT 
					COUNT(r.id) AS boys_count,
					s.day_id
				FROM registrations r 
				LEFT JOIN signups s ON s.registration_id = r.id
				WHERE r.gender='male' and r.deleted_at IS NULL and s.deleted_at IS NULL
				GROUP BY s.day_id
			),
			girls AS (
				SELECT 
					COUNT(r.id) AS girls_count,
					s.day_id
				FROM registrations r 
				LEFT JOIN signups s ON s.registration_id = r.id
				WHERE r.gender='female' and r.deleted_at IS NULL and s.deleted_at IS NULL
				GROUP BY s.day_id
			)
			SELECT 
				d.id AS day_id,
				e.id AS event_id,
				d.capacity,
				d.limit_boys,
				d.limit_girls,
				COALESCE(b.boys_count,0) AS boys_count,
				COALESCE(g.girls_count,0) AS girls_count
			FROM
				days d 
			LEFT JOIN events e ON e.id = d.event_id
			LEFT JOIN boys b ON b.day_id = d.id
			LEFT JOIN girls g ON g.day_id = d.id
			%s
			ORDER BY d.id
		`, where), args...); err != nil {
		return nil, errors.WithStack(err)
	}
	return stats, nil
}

func (repo *PostgresRepo) getStatInternal(ctx context.Context, db sqlx.QueryerContext, eventID int) ([]model.Stat, error) {
	return repo.getStats(ctx, db, "WHERE e.id = $1", eventID)
}

func (repo *PostgresRepo) GetStat(ctx context.Context, eventID int) ([]model.Stat, error) {
	return repo.getStatInternal(ctx, repo.db, eventID)
}

func (repo *PostgresRepo) GetStats(ctx context.Context) ([]model.Stat, error) {
	return repo.getStats(ctx, repo.db, "")
}

func (repo *PostgresRepo) FindOwner(ctx context.Context, username string) (*model.Owner, error) {
	var owner model.Owner
	if err := sqlx.GetContext(ctx, repo.db, &owner, `SELECT * FROM owners WHERE username = $1`, username); err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}
	return &owner, nil
}

func (repo *PostgresRepo) UpdateRegistrations(ctx context.Context, reg *model.Registration) error {
	stmt, err := repo.db.PrepareNamedContext(ctx, `
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

func (repo *PostgresRepo) markRegistrationAsDeleted(ctx context.Context, id int) (*model.Registration, error) {
	var deleted model.Registration
	err := repo.db.GetContext(ctx, &deleted, `
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

func (repo *PostgresRepo) MarkRegistrationAsNotified(ctx context.Context, tx sqlx.QueryerContext, id int) error {
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

func (repo *PostgresRepo) listRegistrations(ctx context.Context, db sqlx.QueryerContext, where string, args ...interface{}) ([]model.ExtendedRegistration, error) {
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
			r.notification_sent_at
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

func (repo *PostgresRepo) ListRegistrations(ctx context.Context) ([]model.ExtendedRegistration, error) {
	return repo.listRegistrations(ctx, repo.db, "r.deleted_at IS NULL AND s.deleted_at IS NULL")
}

func (repo *PostgresRepo) ListEventRegistrations(ctx context.Context, eventID int) ([]model.ExtendedRegistration, error) {
	return repo.listRegistrations(ctx, repo.db, "e.event_id=$1 AND deleted_at IS NULL", eventID)
}

func (repo *PostgresRepo) ListRegistrationsWithoutNotification(ctx context.Context, tx *sqlx.Tx) ([]model.ExtendedRegistration, error) {
	return repo.listRegistrations(ctx, tx, "r.notification_sent_at IS NULL AND r.deleted_at IS NULL AND s.deleted_at IS NULL")
}

func (repo *PostgresRepo) FindRegistrationByID(ctx context.Context, regID int) (*model.ExtendedRegistration, error) {
	regs, err := repo.listRegistrations(ctx, repo.db, "r.id = $1 AND r.deleted_at IS NULL", regID)
	if err != nil {
		return nil, err
	}
	if len(regs) == 0 {
		return nil, errors.WithStack(sql.ErrNoRows)
	}
	return &regs[0], nil
}

func (repo *PostgresRepo) DeleteRegistrationByID(ctx context.Context, regID int) (*model.Registration, error) {
	return repo.markRegistrationAsDeleted(ctx, regID)
}

func (repo *PostgresRepo) FindRegistrationByToken(ctx context.Context, token string) (*model.ExtendedRegistration, error) {
	regs, err := repo.listRegistrations(ctx, repo.db, "r.token = $1 IS NULL", token)
	if err != nil {
		return nil, err
	}
	if len(regs) == 0 {
		return nil, errors.WithStack(sql.ErrNoRows)
	}
	return &regs[0], nil
}

func (repo *PostgresRepo) GeneratePromoCode(ctx context.Context, email string, count int) (token string, err error) {
	err = repo.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		token, err = repo.promoManager.GenerateToken(ctx, tx, email, count)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (repo *PostgresRepo) ValidatePromoCode(ctx context.Context, code string) (int, error) {
	var availableRegistrations int
	err := repo.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		promoCode, err := repo.promoManager.ValidateToken(ctx, tx, code)
		if err != nil {
			return errors.WithStack(err)
		}
		availableRegistrations = promoCode.AvailableRegistrations

		return nil
	})
	if err != nil {
		return 0, err
	}
	return availableRegistrations, nil
}

func (repo *PostgresRepo) WithTxx(ctx context.Context, f func(context.Context, *sqlx.Tx) error) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to begin a transaction.")
	}
	if err := f(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit a transaction.")
	}
	return nil
}
