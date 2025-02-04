package logger

import (
	"context"
	"io"
	"io/ioutil"
	"sync"

	lr "github.com/sirupsen/logrus"
)

const (
	infoLogger string = `Logger:`
	errLogger  string = `%s Logger Error`
	OK         string = "[OK]"
	FAILED     string = "[FAILED]"
)

var (
	once = sync.Once{}
)

type Logger interface {
	SetOptions(opt Options)
	TraceWithContext(ctx context.Context, v ...interface{})
	DebugWithContext(ctx context.Context, v ...interface{})
	InfoWithContext(ctx context.Context, v ...interface{})
	WarnWithContext(ctx context.Context, v ...interface{})
	ErrorWithContext(ctx context.Context, v ...interface{})
	FatalWithContext(ctx context.Context, v ...interface{})
	PanicWithContext(ctx context.Context, v ...interface{})
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Panic(v ...interface{})
}

type logger struct {
	mu       *sync.RWMutex
	log      *lr.Logger
	logEntry *lr.Entry
	opt      Options
}

type Options struct {
	Output        string
	Formatter     string
	Level         string
	ContextFields map[string]string
	DefaultFields map[string]string
	CustomFields  map[string]interface{}
	CustomWriter  io.Writer
	Hook          string
	HookLevel     string
	Hooks         []Hook
}

var (
	DefaultLoggerOption = Options{
		Output:    OutputStdout,
		Formatter: FormatJSON,
		Level:     LevelInfo,
		ContextFields: map[string]string{
			"path":            "x-server-route",
			"request_id":      "x-request-id",
			"method":          "x-request-method",
			"scheme":          "x-request-scheme",
			"user_id":         "x-user-id",
			"client_ip":       "x-forwarded-for",
			"bpm_process_id":  "x-bpm-process-id",
			"bpm_workflow_id": "x-bpm-workflow-id",
			"bpm_instance_id": "x-bpm-instance-id",
			"bpm_job_id":      "x-bpm-job-id",
			"bpm_job_type":    "x-bpm-job-type",
		},
		DefaultFields: map[string]string{
			"name":    "logger",
			"version": "1.0",
		},
		Hook:      "",
		HookLevel: LevelError,
	}
)

var lg *logger

func Init(opt Options) Logger {
	once.Do(func() {
		logrus := lr.New()
		log := logrus.WithFields(lr.Fields{})
		lg = &logger{
			mu:       &sync.RWMutex{},
			log:      logrus,
			logEntry: log,
			opt:      opt,
		}
		lg.log.SetOutput(ioutil.Discard)
		lg.setDefaultOptions()
		lg.applyOptions()
	})
	return lg
}

func (l *logger) setDefaultOptions() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.opt.Output == "" {
		//never put default to discard, error will not be displayed!
		l.opt.Output = DefaultLoggerOption.Output
	}
	if l.opt.Formatter == "" {
		l.opt.Formatter = DefaultLoggerOption.Formatter
	}
	if l.opt.Level == "" {
		l.opt.Level = DefaultLoggerOption.Level
	}
	if l.opt.ContextFields == nil {
		l.opt.ContextFields = DefaultLoggerOption.ContextFields
	}
	if l.opt.DefaultFields == nil {
		l.opt.DefaultFields = DefaultLoggerOption.DefaultFields
	}
	if l.opt.Hook == "" {
		l.opt.Hook = DefaultLoggerOption.Hook
	}
	if l.opt.HookLevel == "" {
		l.opt.HookLevel = DefaultLoggerOption.HookLevel
	}
}

func (l *logger) applyOptions() {
	l.convertAndSetOutput()
	l.convertAndSetFormatter()
	l.convertAndSetLevel()
	l.setDefaultFields()
	l.setHook()
}

func (l *logger) SetOptions(opt Options) {
	l.mu.Lock()
	l.opt = opt
	l.mu.Unlock()
	l.applyOptions()
}
