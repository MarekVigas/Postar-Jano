package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MarekVigas/Postar-Jano/internal/repository"

	"github.com/MarekVigas/Postar-Jano/internal/auth"

	"github.com/MarekVigas/Postar-Jano/internal/api"
	"github.com/MarekVigas/Postar-Jano/internal/db"
	"github.com/MarekVigas/Postar-Jano/internal/mailer"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

const (
	hostEnvKey                = "host"
	portEnvKey                = "port"
	defaultPort               = "5000"
	HTTPServerShutdownTimeout = 5 * time.Second
)

func Run(logger *zap.Logger, fnc func(*tomb.Tomb) error) error {

	// Start catching signals.
	var t tomb.Tomb

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

	// Wait for all threads to drop dead.
	err := t.Wait()

	// In case a special exit code is set, exit.
	if err != nil {
		os.Exit(1)
	}

	return nil
}

func runHTTP(logger *zap.Logger, handler http.Handler) func(t *tomb.Tomb) error {
	return func(t *tomb.Tomb) error {
		host := os.Getenv(hostEnvKey)
		if host == "" {
			host = "0.0.0.0"
		}

		port := os.Getenv(portEnvKey)
		if port == "" {
			port = defaultPort
		}

		logger = logger.With(zap.String("http_port", port))

		// Start the server.
		s := http.Server{
			Addr:    fmt.Sprintf("%v:%s", host, port),
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

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	mailer, err := mailer.NewClient(logger)
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET missing")
	}

	repo := repository.NewPostgresRepo(postgres)

	server := api.New(logger, repo, auth.NewFromDB(repo), mailer, []byte(jwtSecret))

	if err := Run(logger, runHTTP(logger, server)); err != nil {
		logger.Fatal("Failed to run server.", zap.Error(err))
	}
}
