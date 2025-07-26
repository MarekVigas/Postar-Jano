package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PostgresDB struct {
	db *sqlx.DB
}

func NewPostgresDB(db *sql.DB) *PostgresDB {
	return &PostgresDB{db: sqlx.NewDb(db, "postgres")}
}

func (pg *PostgresDB) Ping(ctx context.Context) error {
	if err := pg.db.PingContext(ctx); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (pg *PostgresDB) QueryerContext() sqlx.QueryerContext {
	return pg.db
}

func (pg *PostgresDB) WithTxx(ctx context.Context, f func(context.Context, *sqlx.Tx) error) error {
	tx, err := pg.db.BeginTxx(ctx, nil)
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
