package notification

import (
	"context"
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/logger"

	"gorm.io/gorm"
)

type NotificationDomain struct {
	cfg *configs.AppConfig
	db  *gorm.DB
	log logger.Logger
}

type NotificationDomainHandler interface {
	BeginTx(ctx context.Context) *gorm.DB
	NotificationDomainWriter
}

func NewNotificationDomain(cfg *configs.AppConfig, log logger.Logger, db *gorm.DB) NotificationDomainHandler {
	return &NotificationDomain{
		cfg: cfg,
		db:  db,
		log: log,
	}
}
