package models

import (
	"time"

	"github.com/go-openapi/strfmt"
)

// User model
type User struct {
	ID        strfmt.UUID4 `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Email     string       `gorm:"unique;not null;column:email"`
	Password  string       `gorm:"not null;column:password"`
	CreatedAt time.Time    `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime;column:updated_at"`
	DeletedAt *time.Time   `gorm:"column:deleted_at"`
}

func (user *User) TableName() string {
	return "users"
}
