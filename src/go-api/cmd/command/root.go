package command

import (
	"github.com/MarekVigas/Postar-Jano/internal/api"
	"github.com/MarekVigas/Postar-Jano/internal/auth"
	"github.com/MarekVigas/Postar-Jano/pkg/service"
	"github.com/spf13/cobra"

	"github.com/MarekVigas/Postar-Jano/internal/config"
	"github.com/MarekVigas/Postar-Jano/internal/mailer"
	"github.com/MarekVigas/Postar-Jano/internal/promo"
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

	return service.WithTomb(logger, service.SetupHTTP(&c.Server, logger, handler))
}
