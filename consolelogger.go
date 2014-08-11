package log4go

import (
	"os"
	"time"
)

type consoleAppender struct {
	stdErrAware bool
}

func (a *consoleAppender) Append(msg string, level LogLevel, tstamp time.Time) {
	if a.stdErrAware && level >= LogWarn {
		os.Stderr.Write([]byte(msg))
	} else {
		os.Stdout.Write([]byte(msg))
	}
}

// Create a new console logger that sends all log messages to stdout
func NewConsoleLogger(level LogLevel, timePrefix string) Log4Go {
	appender := consoleAppender{false}
	result := stdLogger{
		&appender,
		level,
		timePrefix,
	}
	return Log4Go(&result)
}

// Create a new conole logger that sends messages less than WARN to stdout and greater that or equal to WARN to stderr
func NewConsoleLoggerWithStderr(level LogLevel, timePrefix string) Log4Go {
	appender := consoleAppender{true}
	result := stdLogger{
		&appender,
		level,
		timePrefix,
	}
	return Log4Go(&result)
}
