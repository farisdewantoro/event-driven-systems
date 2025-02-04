package models

import (
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgtype"
)

type AsynqSendNotificationPayload struct {
	NotificationID   strfmt.UUID4 `json:"notification_id"`
	UserID           strfmt.UUID4 `json:"user_id"`
	NotificationType string       `json:"notification_type"`
}

func (o *AsynqSendNotificationPayload) ToJSON() (*pgtype.JSONB, error) {
	p := &pgtype.JSONB{}
	if err := p.Set(o); err != nil {
		return nil, err
	}
	return p, nil
}
