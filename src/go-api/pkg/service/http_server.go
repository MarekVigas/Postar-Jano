package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MarekVigas/Postar-Jano/internal/config"
	"gopkg.in/tomb.v2"

	"go.uber.org/zap"
)

const (
	HTTPServerShutdownTimeout = 5 * time.Second
)

func SetupHTTP(cfg *config.Server, logger *zap.Logger, handler http.Handler) func(t *tomb.Tomb) error {
	return func(t *tomb.Tomb) error {
		logger = logger.With(zap.Int("http_port", cfg.Port))

		// Start the server.
		s := http.Server{
			Addr:    fmt.Sprintf("%v:%d", cfg.Host, cfg.Port),
			Handler: handler,
		}

		t.Go(func() error {
			logger.Info("Starting the HTTP server...")

			err := s.ListenAndServe()
			if err == http.ErrServerClosed {
				err = nil
			}
			if err != nil {
				logger.Error("HTTP server crashed.", zap.Error(err))
				return err
			}

			logger.Info("HTTP server terminated.")
			return nil
		})

		// Shutdown on tomb dying.
		t.Go(func() error {
			<-t.Dying()

			ctx, cancel := context.WithTimeout(context.Background(), HTTPServerShutdownTimeout)
			defer cancel()

			logger.Info("HTTP server shutdown in progress...", zap.Duration("timeout", HTTPServerShutdownTimeout))
			return s.Shutdown(ctx)
		})

		return nil
	}
}
