package main

import (
	"context"

	"eventdrivensystem/pkg/logger"
	"eventdrivensystem/pkg/logger/mock"

	"go.elastic.co/apm"
)

type CtxKey string

var (
	ContextKey  CtxKey
	ContextKey2 CtxKey
)

func main() {

	// Recommended ways to use this logger
	// Remove default logger code before running this code block
	customLog := logger.Init(logger.Options{
		// Default fields that always included in each log.
		DefaultFields: map[string]string{
			"service.name":    "iwak-tester",
			"service.version": "1.0.0",
		},
		// Only print [INFO] level log or higher.
		// logrus log level can be seen here
		// https://pkg.go.dev/github.com/sirupsen/logrus#:~:text=Logrus%20has%20seven%20logging%20levels,%2C%20Error%2C%20Fatal%20and%20Panic.
		Level: logger.LevelInfo,
		// Custom fields that may be included in each log
		// if found in the passed context.
		// The key on CustomFields will be the key that is printed in the log
		// while the value on CustomFields will be the key that is used to
		// find the value in context.
		// For example:
		// if we call customLog.InfoWithContext(ctx, "info"), with this configuration
		// the function will look for the value for ContextKey and "log.key.2" in
		// the passed context.
		// val := ctx.Value(ContextKey)
		// val := ctx.Value("log.key.2")
		CustomFields: map[string]interface{}{
			"log_key":   ContextKey,
			"log_key_2": ContextKey2,
		},
		// Hook will send log with the same or higher level
		// from HookLevel to specified target. In this case, APM.
		// This will also trigger the logger to do APM trace
		// so span.id, transaction.id, and other APM specific field
		// will be included in each log if any and format the log
		// to follow ElasticAPM standard (ie. key for log level is "log.level" not just "level")
		Hook:      logger.APM,
		HookLevel: logger.LevelError,

		// Hooks are a list of structs that implement the logrus.Hook interface.
		// behind the scenes this will run the command logrus.AddHook()
		// for each hook in the list
		Hooks: []logger.Hook{mock.SampleHook{}},
	})

	tx := apm.DefaultTracer.StartTransaction("name", "type")
	defer tx.End()

	// Init new context with APM specific fields and values inside
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	span, ctx := apm.StartSpan(ctx, "name", "type")
	defer span.End()

	// Add custom fields and values to context
	ctx = context.WithValue(ctx, ContextKey, "taken from ContextKey inside ctx")
	ctx = context.WithValue(ctx, ContextKey2, "taken from log.key.2 inside ctx")

	customLog.Info("info")
	customLog.InfoWithContext(ctx, "info")

	customLog.Error("error")
	customLog.ErrorWithContext(ctx, "error")
}
