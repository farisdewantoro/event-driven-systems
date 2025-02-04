package worker

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/logger"

	"github.com/hibiken/asynq"
)

type WorkerHandler struct {
	cfg *configs.AppConfig
	log logger.Logger
	mux *asynq.ServeMux
}

func NewWorkerHandler(cfg *configs.AppConfig, log logger.Logger, mux *asynq.ServeMux) *WorkerHandler {
	return &WorkerHandler{
		cfg: cfg,
		log: log,
		mux: mux,
	}
}

func (w *WorkerHandler) RegisterHandlers() {
	w.RegisterNotificationHandlers()
}
