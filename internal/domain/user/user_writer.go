package user

import (
	"context"
	models "eventdrivensystem/internal/models/user"
	"eventdrivensystem/pkg/util"

	"gorm.io/gorm"
)

type UserDomainWriter interface {
	CreateUser(ctx context.Context, user *models.User, opts ...util.DbOptions) (*models.User, error)
}

func (u *UserDomain) CreateUser(ctx context.Context, user *models.User, opts ...util.DbOptions) (*models.User, error) {
	return u.createUserSql(ctx, user, opts...)
}

func (u *UserDomain) BeginTx(ctx context.Context) *gorm.DB {
	return u.db.Begin()
}
