package log4go

import (
  "fmt"
  "time"
)

type stdLogger struct {
  level      LogLevel
  appender   appender
  timeFormat string
}

func (l *stdLogger) Log(level LogLevel, format string, a ...interface{}) {
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
