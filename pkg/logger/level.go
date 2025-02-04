package logger

import (
	"fmt"

	lr "github.com/sirupsen/logrus"
)

const (
	LevelTraceMsg   string = "[TRACE]"
	LevelDebugMsg   string = "[DEBUG]"
	LevelInfoMsg    string = "[INFO]"
	LevelWarnMsg    string = "[WARNING]"
	LevelErrorMsg   string = "[ERROR]"
	LevelFatalMsg   string = "[FATAL]"
	LevelPanicMsg   string = "[PANIC]"
	LevelUnknownMsg string = "[UNKNOWN LOG LEVEL]"
)

var (
	ErrUnknownLevel error = fmt.Errorf(`[UNKNOWN LOG LEVEL] [FAILED] Logger Error`)
)

func (l *logger) convertAndSetLevel() {
	l.setLevelLogrus()
}

func (l *logger) setLevelLogrus() {
	lrLevel, lvl := l.getLevel(l.opt.Level)
	//set logrus log level
	l.log.SetLevel(lrLevel)
	l.log.Info(OK, "Logger Level: ", lvl)
}

func (l *logger) getLevel(lvl string) (lr.Level, string) {
	var lrLevel lr.Level
	var lvlString string

	switch lvl {
	case LevelTrace:
		lrLevel = lr.TraceLevel
		lvlString = LevelTraceMsg
	case LevelDebug:
		lrLevel = lr.DebugLevel
		lvlString = LevelDebugMsg
	case LevelInfo:
		lrLevel = lr.InfoLevel
		lvlString = LevelInfoMsg
	case LevelWarn:
		lrLevel = lr.WarnLevel
		lvlString = LevelWarnMsg
	case LevelError:
		lrLevel = lr.ErrorLevel
		lvlString = LevelErrorMsg
	case LevelFatal:
		lrLevel = lr.FatalLevel
		lvlString = LevelFatalMsg
	case LevelPanic:
		lrLevel = lr.PanicLevel
		lvlString = LevelPanicMsg
	default:
		err := ErrUnknownLevel
		l.log.Panic(err)
	}

	return lrLevel, lvlString
}
