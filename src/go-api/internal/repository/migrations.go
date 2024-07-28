package repository

import (
	"github.com/MarekVigas/Postar-Jano/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func RunMigrations(logger *zap.Logger, cfg *config.DB, databaseName string, migrationSrc string) error {
	db, err := cfg.Connect()
	if err != nil {
		return errors.Wrap(err, "failed to connect to DB")
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.WithStack(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationSrc,
		databaseName, driver)
	if err != nil {
		logger.Error("Failed to init migration", zap.Error(err))
		return errors.WithStack(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No migration to run")
			return nil
		}
		logger.Error("Failed to run migration", zap.Error(err))
		return errors.WithStack(err)
	}
	logger.Info("Migrations applied successfully.")
	m.Close()
	return nil
}
