package user

import (
	"context"
	models "eventdrivensystem/internal/models/user"
	"eventdrivensystem/pkg/util"

	"gorm.io/gorm"
)

func (u *UserDomain) createUserSql(ctx context.Context, user *models.User, opts ...util.DbOptions) (*models.User, error) {
	var (
		db  *gorm.DB
		opt util.DbOptions
	)

	if len(opts) > 0 {
		opt = opts[0]
	}

	db = opt.Extract(ctx, u.db)

	return user, db.Create(user).Error
}
