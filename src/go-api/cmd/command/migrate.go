package command

import (
	"github.com/MarekVigas/Postar-Jano/internal/config"
	"github.com/MarekVigas/Postar-Jano/internal/repository"
	"github.com/MarekVigas/Postar-Jano/pkg/service"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "run migrations",
	Run:   migrations,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrations(cmd *cobra.Command, args []string) {
	if err := runMigrations(); err != nil {
		panic(err)
	}
}

func runMigrations() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	return service.WithTomb(logger, func(t *tomb.Tomb) error {
		// End the process
		defer t.Kill(nil)

		c, err := config.LoadAdminSetting()
		if err != nil {
			return errors.Wrap(err, "failed to load config")
		}

		if err := repository.RunMigrations(logger, &c.DB, c.DB.Database, "file://migrations"); err != nil {
			return err
		}

		return nil
	})
}
