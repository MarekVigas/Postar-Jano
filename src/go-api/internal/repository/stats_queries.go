package repository

import (
	"context"
	"fmt"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func getStats(ctx context.Context, db sqlx.QueryerContext, where string, args ...interface{}) ([]model.Stat, error) {
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

func getStatInternal(ctx context.Context, db sqlx.QueryerContext, eventID int) ([]model.Stat, error) {
	return getStats(ctx, db, "WHERE e.id = $1", eventID)
}

func GetStat(ctx context.Context, db sqlx.QueryerContext, eventID int) ([]model.Stat, error) {
	return getStatInternal(ctx, db, eventID)
}

func GetStats(ctx context.Context, db sqlx.QueryerContext) ([]model.Stat, error) {
	return getStats(ctx, db, "")
}
