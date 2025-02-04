package user

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/internal/domain"
	"eventdrivensystem/internal/domain/notification"
	"eventdrivensystem/internal/domain/outbox"
	"eventdrivensystem/internal/domain/user"
	"eventdrivensystem/pkg/logger"
)

type UserUsecase struct {
	cfg *configs.AppConfig
	log logger.Logger

	// domain
	userDomain         user.UserDomainHandler
	outboxDomain       outbox.OutboxDomainHandler
	notificationDomain notification.NotificationDomainHandler
}

type UserUsecaseHandler interface {
	UserUsecaseWriter
}

func NewUserUsecase(
	cfg *configs.AppConfig,
	log logger.Logger,
	dom *domain.Domain,
) UserUsecaseHandler {
	return &UserUsecase{
		cfg:                cfg,
		log:                log,
		userDomain:         dom.User,
		outboxDomain:       dom.Outbox,
		notificationDomain: dom.Notification,
	}
}
