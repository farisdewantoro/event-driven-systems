package logger

import (
	"fmt"

	lr "github.com/sirupsen/logrus"
)

const (
	FormatJSON string = "json"
	FormatText string = "text"

	formatJSON string = "[JSON]"
	formatText string = "[TEXT]"

	loggerFormat string = "Logger Formatter: "
)

var (
	ErrUnknownFormat = fmt.Errorf(`[UNKNOWN LOG FORMAT] [FAILED] Logger Error`)
	jsonFormatter    = &lr.JSONFormatter{}
	textFormatter    = &lr.TextFormatter{}
)

func (l *logger) convertAndSetFormatter() {
	switch l.opt.Formatter {
	case FormatText:
		l.log.SetFormatter(textFormatter)
		l.log.Info(OK, loggerFormat, formatText)
	case FormatJSON:
		if l.opt.Hook == APM {
			jsonFormatter.FieldMap = lr.FieldMap{
				lr.FieldKeyLevel: "log.level",
				lr.FieldKeyTime:  "@timestamp",
				lr.FieldKeyMsg:   "message",
			}
		}

		l.log.SetFormatter(jsonFormatter)
		l.log.Info(OK, loggerFormat, formatJSON)
	default:
		l.log.Panic(ErrUnknownFormat)
	}
}
