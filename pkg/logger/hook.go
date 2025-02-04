package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.elastic.co/apm/module/apmlogrus"
)

const (
	APM  string = "apm"
	OTEL string = "otel"
)

// Hook is an alias of logrus.Hook...
// A hook to be fired when logging on the logging levels returned from
// `Levels()` on your implementation of the interface. Note that this is not
// fired in a goroutine or a channel with workers, you should handle such
// functionality yourself if your call is non-blocking and you don't wish for
// the logging calls for levels returned from `Levels()` to block.
type Hook logrus.Hook

func (l *logger) setHook() {
	level, lvl := l.getLevel(l.opt.HookLevel)

	switch l.opt.Hook {
	case APM:
		l.log.AddHook(
			&apmlogrus.Hook{
				LogLevels: []logrus.Level{level},
			})
		l.log.Info(OK, "Logger Hook APM Level: ", lvl)
	case OTEL:
		l.log.AddHook(
			otellogrus.NewHook(otellogrus.WithLevels(level)),
		)
		l.log.Info(OK, "Logger Hook OTEL Level: ", lvl)
	}

	for i, hook := range l.opt.Hooks {
		l.log.AddHook(hook)
		l.log.Info(OK, fmt.Sprintf("Add Additional Hook [%d], level:%v ", i+1, hook.Levels()))
	}
}
