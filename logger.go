package log5go

import (
	"fmt"
	"time"
)

// Inner type of all loggers
type logger struct {
	level      LogLevel
	appender   Appender
	timeFormat string
}

// Log a message at the given log level
func (l *logger) Log(level LogLevel, format string, a ...interface{}) {
	tstamp := time.Now()
	if level >= l.level {
		msg := fmt.Sprintf(format, a...)
		if l.timeFormat != "" {
			timePrefix := tstamp.Format(l.timeFormat)
			msg = fmt.Sprintf("%s %s: %s\n", timePrefix, levelMap[level], msg)
		} else {
			msg = fmt.Sprintf("%s: %s\n", levelMap[level], msg)
		}

		l.appender.Append(msg, level, tstamp)
	}
}

func (l *logger) Trace(format string, a ...interface{}) {
	l.Log(LogTrace, format, a...)
}

func (l *logger) Debug(format string, a ...interface{}) {
	l.Log(LogDebug, format, a...)
}

func (l *logger) Info(format string, a ...interface{}) {
	l.Log(LogInfo, format, a...)
}

func (l *logger) Warn(format string, a ...interface{}) {
	l.Log(LogWarn, format, a...)
}

func (l *logger) Error(format string, a ...interface{}) {
	l.Log(LogError, format, a...)
}

func (l *logger) Fatal(format string, a ...interface{}) {
	l.Log(LogFatal, format, a...)
}

func (l *logger) GetLogLevel() LogLevel {
	return l.level
}

func (l *logger) SetLogLevel(level LogLevel) {
	l.level = level
}
