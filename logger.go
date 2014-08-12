package log4go

import (
	"fmt"
	"time"
)

const (
	TF_NCSA  = "02/Jan/2006:15:04:05 -0700"
	TF_GoStd = "2006/01/02 15:04:05"
)

type Log4Go interface {
	Log(level LogLevel, format string, a ...interface{})
	Trace(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
	GetLogLevel() LogLevel
	SetLogLevel(level LogLevel)
}

type stdLogger struct {
	appender   appender
	level      LogLevel
	timePrefix string
}

func (l *stdLogger) Log(level LogLevel, format string, a ...interface{}) {
	tstamp := time.Now()
	if level >= l.level {
		msg := fmt.Sprintf(format, a...)
		if l.timePrefix != "" {
			timePrefix := tstamp.Format(l.timePrefix)
			msg = fmt.Sprintf("%s %s: %s\n", timePrefix, levelMap[level], msg)
		} else {
			msg = fmt.Sprintf("%s: %s\n", levelMap[level], msg)
		}

		l.appender.Append(msg, level, tstamp)
	}
}

func (l *stdLogger) Trace(format string, a ...interface{}) {
	l.Log(LogTrace, format, a...)
}

func (l *stdLogger) Debug(format string, a ...interface{}) {
	l.Log(LogDebug, format, a...)
}

func (l *stdLogger) Info(format string, a ...interface{}) {
	l.Log(LogInfo, format, a...)
}

func (l *stdLogger) Warn(format string, a ...interface{}) {
	l.Log(LogWarn, format, a...)
}

func (l *stdLogger) Error(format string, a ...interface{}) {
	l.Log(LogError, format, a...)
}

func (l *stdLogger) Fatal(format string, a ...interface{}) {
	l.Log(LogFatal, format, a...)
}

func (l *stdLogger) GetLogLevel() LogLevel {
	return l.level
}

func (l *stdLogger) SetLogLevel(level LogLevel) {
	l.level = level
}
