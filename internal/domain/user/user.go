package user

import (
	"context"
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/logger"

	"gorm.io/gorm"
)

type UserDomain struct {
	cfg *configs.AppConfig
	db  *gorm.DB
	log logger.Logger
}

type UserDomainHandler interface {
	BeginTx(ctx context.Context) *gorm.DB
	UserDomainWriter
}

func NewUserDomain(cfg *configs.AppConfig, log logger.Logger, db *gorm.DB) UserDomainHandler {
	return &UserDomain{
		cfg: cfg,
		db:  db,
		log: log,
	}
}
