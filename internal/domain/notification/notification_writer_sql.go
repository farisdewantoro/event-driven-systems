package notification

import (
	"context"
	models "eventdrivensystem/internal/models/notification"
	"eventdrivensystem/pkg/util"

	"gorm.io/gorm"
)

func (u *NotificationDomain) createNotificationSql(ctx context.Context, p *models.Notification, opts ...util.DbOptions) (*models.Notification, error) {
	var (
		db  *gorm.DB
		opt util.DbOptions
	)

	if len(opts) > 0 {
		opt = opts[0]
	}

	db = opt.Extract(ctx, u.db)

	return p, db.Create(p).Error
}
