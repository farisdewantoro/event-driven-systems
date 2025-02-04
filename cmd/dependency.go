package cmd

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/databases"
	"eventdrivensystem/pkg/logger"
	"log"
	"runtime"
	"sync"

	goValidator "github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

var (
	appDependency         *AppDependency
	appDependencySyncOnce sync.Once
)

type AppDependency struct {
	db        *gorm.DB
	cfg       *configs.AppConfig
	queue     *asynq.Client
	log       logger.Logger
	validator *goValidator.Validate
}

func GetAppDependency() *AppDependency {
	appDependencySyncOnce.Do(func() {
		appDependency = NewAppDependency()
	})
	return appDependency
}

func NewAppDependency() *AppDependency {
	cfg := configs.Get()
	db, err := databases.NewSqlDb(cfg)
	lgOptions := logger.Options{
		Output:    logger.OutputStdout,
		Formatter: logger.FormatJSON,
		Level:     logger.LevelInfo,
		DefaultFields: map[string]string{
			"app.name":    cfg.Meta.Name,
			"app.runtime": runtime.Version(),
		},
	}
	lg := logger.Init(lgOptions)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Redis.Address})

	return &AppDependency{
		db:        db,
		cfg:       cfg,
		queue:     asynqClient,
		log:       lg,
		validator: goValidator.New(),
	}
}
