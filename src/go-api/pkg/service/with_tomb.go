package service

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

func WithTomb(logger *zap.Logger, fnc func(*tomb.Tomb) error) error {
	var t tomb.Tomb

	// Setup signal handlers.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Handle signals.
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

	if err := fnc(&t); err != nil {
		t.Kill(err)
	}

	return t.Wait()
}
