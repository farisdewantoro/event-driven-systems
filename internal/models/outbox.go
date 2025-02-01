package models

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgtype"
)

const (
	OutboxStatusPending    string = "PENDING"
	OutboxStatusProcessing string = "PROCESSING"
	OutboxStatusSent       string = "SENT"
	OutboxStatusFailed     string = "FAILED"
	OutboxStatusRetrying   string = "RETRYING"
)

type Outbox struct {
	ID              strfmt.UUID4 `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Payload         pgtype.JSONB `json:"payload" gorm:"type:jsonb;default:'{}';not null"`
	MessageType     string       `json:"message_type" gorm:"not null"`
	Status          string       `json:"status" gorm:"column:status;not null"`
	CreatedAt       time.Time    `json:"created_at" gorm:"column:created_at"`
	Attempt         int64        `json:"attempt" gorm:"column:attempt"`
	DestinationType string       `json:"destination_type" gorm:"column:destination_type;not null"`
	SentAt          *time.Time   `json:"sent_at,omitempty" gorm:"column:sent_at"`
	ErrorMessage    *string      `json:"error_message,omitempty" gorm:"column:error_message"`
	ExecuteAt       time.Time    `json:"execute_at" gorm:"column:execute_at;not null"`
}

func (Outbox) TableName() string {
	return "outbox"
}
