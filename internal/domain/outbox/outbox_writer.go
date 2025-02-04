package outbox

import (
	"context"
	models "eventdrivensystem/internal/models/outbox"
	"eventdrivensystem/pkg/util"
)

type OutboxDomainWriter interface {
	CreateOutbox(ctx context.Context, Outbox *models.Outbox, opts ...util.DbOptions) error
}

func (u *OutboxDomain) CreateOutbox(ctx context.Context, outbox *models.Outbox, opts ...util.DbOptions) error {
	outbox.Status = models.OutboxStatusPending
	return u.createOutboxSql(ctx, outbox, opts...)
}
