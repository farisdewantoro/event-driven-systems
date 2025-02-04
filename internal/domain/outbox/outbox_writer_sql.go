package outbox

import (
	"context"
	models "eventdrivensystem/internal/models/outbox"
	"eventdrivensystem/pkg/util"

	"gorm.io/gorm"
)

func (u *OutboxDomain) createOutboxSql(ctx context.Context, outbox *models.Outbox, opts ...util.DbOptions) error {
	var (
		db  *gorm.DB
		opt util.DbOptions
	)

	if len(opts) > 0 {
		opt = opts[0]
	}

	db = opt.Extract(ctx, u.db)

	return db.Create(outbox).Error
}
