package databases

import (
	"fmt"
	"loanservice/configs"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Required for postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Required for file source
)

func NewSqlDb(cfg *configs.AppConfig) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.SQL.DSN,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return
	}

	db.Use(tracing.NewPlugin())

	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	// Create the connection pool
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(30)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(620 * time.Second)

	log.Info("Database connected successfully", cfg.SQL.DSN)
	return db, nil
}

func NewMigrate(cfg *configs.AppConfig) (*migrate.Migrate, error) {
	migrationsPath := "file://./docs/db/migrations"

	// Create the migrator instance
	migrator, err := migrate.New(migrationsPath, cfg.SQL.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator instance: %w", err)
	}

	return migrator, nil

}
