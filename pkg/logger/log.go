package logger

import (
	"context"

	lr "github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"
	"go.opentelemetry.io/otel/trace"
)

const (
	LevelTrace string = "trace"
	LevelDebug string = "debug"
	LevelInfo  string = "info"
	LevelWarn  string = "warn"
	LevelError string = "error"
	LevelFatal string = "fatal"
	LevelPanic string = "panic"

	// FieldKeyTraceID is the field key for the trace ID.
	FieldKeyTraceID = "trace.id"

	// FieldKeyTransactionID is the field key for the transaction ID.
	FieldKeyTransactionID = "transaction.id"

	// FieldKeySpanID is the field key for the span ID.
	FieldKeySpanID = "span.id"
)

func (l *logger) parseContextFields(ctx context.Context) *lr.Entry {
	doLog := l.logEntry
	if ctx != nil {
		for k, v := range l.opt.ContextFields {
			if val := ctx.Value(v); val != nil {
				doLog = doLog.WithField(k, val)
			}
		}

		if len(l.opt.CustomFields) != 0 {
			for k, v := range l.opt.CustomFields {
				if val := ctx.Value(v); val != nil {
					doLog = doLog.WithField(k, val)
				}
			}
		}

		switch l.opt.Hook {
		case APM:
			doLog = doLog.WithFields(apmlogrus.TraceContext(ctx))
		case OTEL:
			span := trace.SpanFromContext(ctx)
			traceID := span.SpanContext().TraceID().String()
			spanID := span.SpanContext().SpanID().String()

			fields := lr.Fields{
				FieldKeyTraceID: traceID,
				FieldKeySpanID:  spanID,
			}

			doLog = doLog.WithFields(fields)
		}
	}

	return doLog
}

func (l *logger) TraceWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Trace(v...)
}

func (l *logger) Trace(v ...interface{}) {
	l.TraceWithContext(context.TODO(), v...)
}

func (l *logger) DebugWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Debug(v...)
}

func (l *logger) Debug(v ...interface{}) {
	l.DebugWithContext(context.TODO(), v...)
}

func (l *logger) InfoWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Info(v...)
}

func (l *logger) Info(v ...interface{}) {
	l.InfoWithContext(context.TODO(), v...)
}

func (l *logger) WarnWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Warn(v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.WarnWithContext(context.TODO(), v...)
}

func (l *logger) ErrorWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Error(v...)
}

func (l *logger) Error(v ...interface{}) {
	l.ErrorWithContext(context.TODO(), v...)
}

func (l *logger) FatalWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Fatal(v...)
}

func (l *logger) Fatal(v ...interface{}) {
	l.FatalWithContext(context.TODO(), v...)
}

func (l *logger) PanicWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Panic(v...)
}

func (l *logger) Panic(v ...interface{}) {
	l.PanicWithContext(context.TODO(), v...)
}
