package status

import (
	"context"
	"github.com/MarekVigas/Postar-Jano/internal/repository"
	"github.com/MarekVigas/Postar-Jano/internal/resources"
	"github.com/MarekVigas/Postar-Jano/pkg/logger"
	"go.uber.org/zap"
)

type Checker struct {
	postgresDB *repository.PostgresDB
}

func NewChecker(postgresDB *repository.PostgresDB) *Checker {
	return &Checker{postgresDB: postgresDB}
}

func (c *Checker) Ping(ctx context.Context) (resources.StatusResponse, bool) {
	if err := c.postgresDB.Ping(ctx); err != nil {
		logger.FromCtx(ctx).Error("Failed to ping DB", zap.Error(err))
		return resources.NewErrorStatusResponse(), false
	}
	return resources.NewOkStatusResponse(), true
}
