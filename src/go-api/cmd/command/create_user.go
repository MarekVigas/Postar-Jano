package command

import (
	"context"
	"github.com/MarekVigas/Postar-Jano/internal/repository"

	"github.com/MarekVigas/Postar-Jano/internal/config"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/pkg/service"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

var (
	username string
	password string
)

// addUserCmd adds a new user
var addUserCmd = &cobra.Command{
	Use:   "add-user",
	Short: "Adds a new user",
	Run:   addUser,
}

func init() {
	addUserCmd.PersistentFlags().StringVar(&username, "username", "user@example.com", "username for the create user")
	addUserCmd.PersistentFlags().StringVar(&password, "password", "", "user password")

	rootCmd.AddCommand(addUserCmd)
}

func addUser(cmd *cobra.Command, args []string) {
	if err := runAddUser(); err != nil {
		panic(err)
	}
}

func runAddUser() error {
	if password == "" {
		return errors.New("password must be set")
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	return service.WithTomb(logger, func(t *tomb.Tomb) error {
		c, err := config.LoadAdminSetting()
		if err != nil {
			return errors.Wrap(err, "failed to load config")
		}

		db, err := c.DB.Connect()
		if err != nil {
			return errors.Wrap(err, "failed to connect to DB")
		}
		defer db.Close()

		ctx := t.Context(context.Background())
		_, err = repository.CreateOwner(ctx, sqlx.NewDb(db, "postgres"), model.Owner{
			Username: username,
		}, password)

		if err != nil {
			logger.Error("Failed to create user.", zap.Error(err))
			return err
		}
		logger.Info("User created successfully.")

		// End the process.
		t.Kill(nil)
		return nil
	})
}
