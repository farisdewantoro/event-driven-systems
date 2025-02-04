package worker

import (
	"context"
	"encoding/json"
	models "eventdrivensystem/internal/models/asynq"

	"github.com/hibiken/asynq"
)

func (w *WorkerHandler) RegisterNotificationHandlers() {
	w.mux.HandleFunc(models.AsynqTaskSendEmailNotification, w.handleSendEmailNotification)
}

func (w *WorkerHandler) handleSendEmailNotification(ctx context.Context, task *asynq.Task) error {
	var (
		param = models.AsynqSendNotificationPayload{}
	)

	if err := json.Unmarshal(task.Payload(), &param); err != nil {
		return err
	}

	w.log.InfoWithContext(ctx, "Email verification request sent for UserID: "+param.UserID)

	return nil
}
