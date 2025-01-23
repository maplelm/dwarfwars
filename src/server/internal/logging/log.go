package logging

import (
	"fmt"
	zlog "github.com/rs/zerolog/log"
)

type LogLevel int

const (
	None LogLevel = iota
	Trace
	Debug
	System
	Warning
	Error
)

var CurrentLoggingLevel LogLevel = Trace

func Printf(l LogLevel, format string, args ...any) {
	if l >= CurrentLoggingLevel {
		switch l {
		case Trace:
			zlog.Trace().Str("level", "Trace").Msg(fmt.Sprintf(format, args...))
		case Debug:
			zlog.Debug().Str("level", "Debug").Msg(fmt.Sprintf(format, args...))
		case System:
		case Warning:
		case Error:
		default:
		}
	}
}

func Print(l LogLevel, msg string) {
	Printf(l, "%s", msg)
}
