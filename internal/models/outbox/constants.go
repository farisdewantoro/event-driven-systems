package models

const (
	OutboxStatusPending    string = "PENDING"
	OutboxStatusProcessing string = "PROCESSING"
	OutboxStatusSent       string = "SENT"
	OutboxStatusFailed     string = "FAILED"
	OutboxStatusRetrying   string = "RETRYING"

	OutboxDestinationTypeKafka    string = "KAFKA"
	OutboxDestinationTypeRabbitmq string = "RABBITMQ"
	OutboxDestinationTypeAsynq    string = "ASYNQ"
)
