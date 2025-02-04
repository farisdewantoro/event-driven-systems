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
	Meta      Meta
	ApiServer ApiServer
	SQL       SQL
	Redis     Redis
	Outbox    Outbox
	AsyncQ    AsyncQ
}

type Meta struct {
	Name string
}

type ApiServer struct {
	Host string `validate:"required"`
	Port int    `validate:"required"`
}

type SQL struct {
	DSN string `validate:"required"`
}

type Redis struct {
	Address string `validate:"required"`
}

type Outbox struct {
	MaxRetries           int
	MaxConcurrency       int
	MaxBatchSize         int
	DurationIntervalInMs int
}

type AsyncQ struct {
	MaxRetries              int
	BasedServiceConsumerURL string `validate:"required"`
	MonitoringHost          string `validate:"required"`
	MonitoringPort          int    `validate:"required"`
}

func Get() *AppConfig {

	if cfg == nil {
		Load()
	}

	return cfg
}

func Load() {
	onceCfg.Do(func() {
		v := viper.New()
		v.AddConfigPath(".")
		v.SetConfigType("yaml")
		v.SetConfigName("config")

		err := v.ReadInConfig()
		if err != nil {
			log.Fatalf("error reading config file, %v", err)
		}

		c := &AppConfig{}
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
