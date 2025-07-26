package logger

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

type ctxKeyLogger struct{}

func FromCtx(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxKeyLogger{}).(*zap.Logger); ok {
		return logger
	}
	panic("no logger found in context")
}

func UnexpectedError(ctx context.Context, err error) *zap.Logger {
	logger := FromCtx(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return zap.NewNop()
	default:
		return logger
	}
}
