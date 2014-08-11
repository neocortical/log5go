package log4go

import (
  "time"
)

type appender interface {
  Append(msg string, level LogLevel, tstamp time.Time)
}
