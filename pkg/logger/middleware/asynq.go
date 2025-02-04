package middleware

import (
	"context"
	"eventdrivensystem/pkg/logger"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

func LoggingMiddlewareAsynq(lg logger.Logger) func(handler asynq.Handler) asynq.Handler {
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {

			var (
				start  = time.Now()
				taskID = t.ResultWriter().TaskID()
			)

			lg.Info(fmt.Sprintf("Start process task %s(%s)", t.Type(), taskID))
			lg.Info("payload :", string(t.Payload()))
			err := h.ProcessTask(ctx, t)
			if err != nil {

				return err
			}
			lg.Info(fmt.Sprintf("Finished processing %s(%s), duration: %v", t.Type(), taskID, time.Since(start)))
			return nil
		})
	}
}
