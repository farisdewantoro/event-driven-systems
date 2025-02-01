package configs

import (
	"sync"

	goValidator "github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

var (
	cfg     *AppConfig
	onceCfg sync.Once
)

type AppConfig struct {
	ApiServer ApiServer
	SQL       SQL
	Redis     Redis
	Outbox    Outbox
}

// API Server configuration
type ApiServer struct {
	Name string
	Host string
	Port int
}

type SQL struct {
	DSN string
}

type Redis struct {
	Address string
}

type Outbox struct {
	MaxRetries           int
	MaxConcurrency       int
	MaxBatchSize         int
	DurationIntervalInMs int
}

func Get() *AppConfig {
	return cfg
}

func Load() {
	onceCfg.Do(func() {
		v := viper.New()
		v.AutomaticEnv()

		v.AddConfigPath(".")
		v.SetConfigType("json")
		v.SetConfigName("config")

		err := v.ReadInConfig()
		if err != nil {
			log.Fatalf("error reading config file, %v", err)
		}

		c := &AppConfig{
			ApiServer: ApiServer{
				Name: v.GetString("api_server.name"),
				Host: v.GetString("api_server.host"),
				Port: v.GetInt("api_server.port"),
			},
			SQL: SQL{
				DSN: v.GetString("sql.dsn"),
			},
			Redis: Redis{
				Address: v.GetString("redis.address"),
			},
			Outbox: Outbox{
				MaxRetries:           v.GetInt("outbox.max_retries"),
				MaxConcurrency:       v.GetInt("outbox.max_concurrency"),
				MaxBatchSize:         v.GetInt("outbox.max_batch_size"),
				DurationIntervalInMs: v.GetInt("outbox.duration_interval_in_ms"),
			},
		}
		err = v.Unmarshal(&c)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}

		// Perform validation
		validate := goValidator.New()
		if err := validate.Struct(c); err != nil {
			if validationErrors, ok := err.(goValidator.ValidationErrors); ok {
				for _, ve := range validationErrors {
					log.Printf("Validation error for field '%s': %s", ve.StructNamespace(), ve.Tag())
				}
			} else {
				log.Fatalf("Validation error: %v", err)
			}
			log.Fatalf("Invalid config")
		}

		cfg = c

	})
}
