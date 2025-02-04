package notification

import (
	"context"
	models "eventdrivensystem/internal/models/notification"
	"eventdrivensystem/pkg/util"

	"gorm.io/gorm"
)

type NotificationDomainWriter interface {
	CreateNotification(ctx context.Context, p *models.Notification, opts ...util.DbOptions) (*models.Notification, error)
}

func (u *NotificationDomain) CreateNotification(ctx context.Context, p *models.Notification, opts ...util.DbOptions) (*models.Notification, error) {
	return u.createNotificationSql(ctx, p, opts...)
}

func (u *NotificationDomain) BeginTx(ctx context.Context) *gorm.DB {
	return u.db.Begin()
}
