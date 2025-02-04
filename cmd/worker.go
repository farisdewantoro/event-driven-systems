package cmd

import (
	"context"
	"eventdrivensystem/internal/handler/worker"
	"eventdrivensystem/pkg/logger/middleware"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var asynqWorkerCmd = &cobra.Command{
	Use:   "asynq-worker",
	Short: "Start the worker service",
	Run: func(cmd *cobra.Command, args []string) {
		go StartWorkerService()
		StartWorker()
	},
}

func StartWorkerService() {
	var srv http.Server

	dp := GetAppDependency()

	idleConnection := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		dp.log.Info("[WORKER-API] is shutting down")
		if err := srv.Shutdown(context.Background()); err != nil {
			dp.log.Info("[WORKER-API] Fail shutting down:", err)
		}
		close(idleConnection)
	}()

	e := echo.New()
	e.Use(echoMiddleware.Recover())
	asynqMon := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitoring/tasks",
		RedisConnOpt: asynq.RedisClientOpt{Addr: dp.cfg.Redis.Address},
	})

	e.Any("/monitoring/tasks/*", echo.WrapHandler(asynqMon))

	go func() {
		address := fmt.Sprintf("%s:%d", dp.cfg.AsyncQ.MonitoringHost, dp.cfg.AsyncQ.MonitoringPort)
		if err := e.Start(address); err != nil {
			dp.log.Info("shutting down the worker server -> ", err)
		}
	}()
}

func StartWorker() {
	dp := GetAppDependency()

	// Asynq Worker Setup
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: dp.cfg.Redis.Address},
		asynq.Config{Concurrency: 0},
	)

	mux := asynq.NewServeMux()
	mux.Use(middleware.LoggingMiddlewareAsynq(dp.log))
	worker.NewWorkerHandler(dp.cfg, dp.log, mux).RegisterHandlers()

	dp.log.Info("Worker started, waiting for tasks...")
	if err := server.Run(mux); err != nil {
		dp.log.Error("Could not start worker: %v", err)
		return
	}
}
