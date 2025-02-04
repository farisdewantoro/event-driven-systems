package outbox

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/logger"

	"gorm.io/gorm"
)

type OutboxDomain struct {
	cfg *configs.AppConfig
	db  *gorm.DB
	log logger.Logger
}

type OutboxDomainHandler interface {
	OutboxDomainWriter
}

func NewOutboxDomain(cfg *configs.AppConfig, log logger.Logger, db *gorm.DB) OutboxDomainHandler {
	return &OutboxDomain{
		cfg: cfg,
		db:  db,
		log: log,
	}
}
