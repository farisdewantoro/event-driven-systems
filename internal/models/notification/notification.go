package notification

import (
	"time"

	"github.com/go-openapi/strfmt"
)

// Notification model
type Notification struct {
	ID        strfmt.UUID4 `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	UserID    strfmt.UUID4 `gorm:"type:uuid;not null;column:user_id"`
	Type      string       `gorm:"not null;column:type"`
	Message   string       `gorm:"not null;column:message"`
	Status    string       `gorm:"not null;column:status"`
	CreatedAt time.Time    `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime;column:updated_at"`
	DeletedAt *time.Time   `gorm:"column:deleted_at"`
}

func (n *Notification) TableName() string {
	return "notifications"
}
