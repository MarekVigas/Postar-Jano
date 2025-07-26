package logger

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

func ContextLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create request-scoped logger
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = "req-" + time.Now().Format("20060102150405")
			}

			logger := logger.With(
				zap.String("request_id", reqID),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path))

			// Store logger into context
			ctx := context.WithValue(c.Request().Context(), ctxKeyLogger{}, logger)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
