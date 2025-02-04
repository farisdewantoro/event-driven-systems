package cmd

import (
	"errors"
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/databases"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrateUpCmd = &cobra.Command{
	Use:   "db-migrate-up",
	Short: "Runs the db migrations",
	Run: func(cmd *cobra.Command, args []string) {
		MigrateUp()
	},
}

func MigrateUp() {
	cfg := configs.Get()
	m, err := databases.NewMigrate(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("no migrations to run -> %v", err)
		} else {
			log.Fatalf("failed to run migrations: %v", err)
		}
		return
	}

	log.Println("migrations run successfully")

}
