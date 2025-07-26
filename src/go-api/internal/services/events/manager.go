package events

import (
	"context"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/repository"
	"github.com/MarekVigas/Postar-Jano/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Manager struct {
	postgresDB *repository.PostgresDB
}

func NewManager(db *repository.PostgresDB) *Manager {
	return &Manager{postgresDB: db}
}

func (m *Manager) GetAll(ctx context.Context) ([]model.Event, error) {
	var events []model.Event
	err := m.postgresDB.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var err error
		events, err = repository.ListEvents(ctx, tx)
		if err != nil {
			logger.FromCtx(ctx).Error("Failed to list events.", zap.Error(err))
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (m *Manager) GetByID(ctx context.Context, eventID int) (*model.Event, error) {
	var event *model.Event
	err := m.postgresDB.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var err error
		event, err = repository.FindEvent(ctx, tx, eventID)
		if err != nil {
			logger.UnexpectedError(ctx, err).Error("Failed to find event.", zap.Error(err), zap.Int("event_id", eventID))
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (m *Manager) GetAllStats(ctx context.Context) ([]model.Stat, error) {
	stats, err := repository.GetStats(ctx, m.postgresDB.QueryerContext())
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to get stats", zap.Error(err))
		return nil, err
	}
	return stats, nil
}

func (m *Manager) GetStatById(ctx context.Context, eventId int) ([]model.Stat, error) {
	stats, err := repository.GetStat(ctx, m.postgresDB.QueryerContext(), eventId)
	if err != nil {
		logger.UnexpectedError(ctx, err).Error("Failed to get stats", zap.Int("event_id", eventId))
		return nil, err
	}
	return stats, nil
}
