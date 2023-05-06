package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MarekVigas/Postar-Jano/internal/api"
	"github.com/MarekVigas/Postar-Jano/internal/auth"
	"github.com/MarekVigas/Postar-Jano/internal/config"
	"github.com/MarekVigas/Postar-Jano/internal/mailer"
	"github.com/MarekVigas/Postar-Jano/internal/promo"
	"github.com/MarekVigas/Postar-Jano/internal/repository"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

const (
	HTTPServerShutdownTimeout = 5 * time.Second
)

func Run(logger *zap.Logger, fnc func(*tomb.Tomb) error) error {
	var t tomb.Tomb

	// Setup signal handlers.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Kill the tomb on signal.
	t.Go(func() error {
		select {
		case sig := <-sigCh:
			logger.Info("Signal received, terminating...", zap.String("signal", sig.String()))

			// Kill the tomb.
			t.Kill(nil)

			// Get killed on the next signal.
			signal.Stop(sigCh)

		case <-t.Dying():
		}
		return nil
	})

	// Call business logic.
	if err := fnc(&t); err != nil {
		t.Kill(err)
	}

	return t.Wait()
}

func setupHTTP(cfg *config.Server, logger *zap.Logger, handler http.Handler) func(t *tomb.Tomb) error {
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

func runMain() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return errors.Wrap(err, "failed to setup logger")
	}

	c, err := config.LoadAdminSetting()
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	postgres, err := c.DB.Connect()
	if err != nil {
		return errors.Wrap(err, "failed to connect to DB")
	}
	defer postgres.Close()

	mailer, err := mailer.NewClient(&c.Mailer, logger)
	if err != nil {
		return errors.Wrap(err, "failed to setup mailer")
	}

	var promoGenerator repository.PromoManager
	if c.Promo.Simple {
		promoGenerator = promo.NewSimpleGenerator(logger)
	} else {
		promoGenerator = promo.NewJWTGenerator(logger, c.Promo.Secret, c.Promo.ActivationDate, c.Promo.ExpirationDate)
	}

	repo := repository.NewPostgresRepo(postgres, promoGenerator)

	handler := api.New(logger, repo, auth.NewFromDB(repo), mailer, c.JWTSecret)

	return Run(logger, setupHTTP(&c.Server, logger, handler))
}

func main() {
	if err := runMain(); err != nil {
		panic(err)
	}
}
