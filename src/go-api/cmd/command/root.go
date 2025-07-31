package command

import (
	"github.com/MarekVigas/Postar-Jano/internal/api"
	"github.com/MarekVigas/Postar-Jano/internal/services/auth"
	"github.com/MarekVigas/Postar-Jano/internal/services/events"
	"github.com/MarekVigas/Postar-Jano/internal/services/mailer"
	"github.com/MarekVigas/Postar-Jano/internal/services/promo"
	"github.com/MarekVigas/Postar-Jano/internal/services/registration"
	"github.com/MarekVigas/Postar-Jano/internal/services/status"
	"github.com/MarekVigas/Postar-Jano/pkg/service"
	"github.com/spf13/cobra"

	"github.com/MarekVigas/Postar-Jano/internal/config"
	"github.com/MarekVigas/Postar-Jano/internal/repository"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	rootCmd = &cobra.Command{
		Use:   "registrations_api",
		Short: "HTTP API for registrations",
		Run:   runRootCmd,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func runRootCmd(cmd *cobra.Command, args []string) {
	if err := runMain(); err != nil {
		panic(err)
	}
}

func setupPromoRegistry(logger *zap.Logger, postgresDB *repository.PostgresDB, mailer *mailer.Client, config config.Promo) *promo.Registry {
	var promoGenerator promo.Generator
	if config.Simple {
		promoGenerator = promo.NewSimpleGenerator(logger)
	} else {
		promoGenerator = promo.NewJWTGenerator(logger, config.Secret, config.ActivationDate, config.ExpirationDate)
	}
	return promo.NewRegistry(postgresDB, promoGenerator, mailer)
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

	mailerClient, err := mailer.NewClient(&c.Mailer)
	if err != nil {
		return errors.Wrap(err, "failed to setup mailerClient")
	}

	postgresDB := repository.NewPostgresDB(postgres)

	promoRegistry := setupPromoRegistry(logger, postgresDB, mailerClient, c.Promo)

	handler := api.New(
		logger,
		auth.NewFromDB(postgresDB, c.JWTSecret),
		events.NewManager(postgresDB),
		promoRegistry,
		registration.NewManager(postgresDB, promoRegistry, mailerClient),
		status.NewChecker(postgresDB),
	)

	return service.WithTomb(logger, service.SetupHTTP(&c.Server, logger, handler))
}
