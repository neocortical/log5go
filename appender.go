package log4go

import (
	"time"
)

// Appender defines the interface responsible for writing a processed message to
// a destination. Currently, destinations are console and file. Additional
// appender types are planned.
type Appender interface {
	Append(msg string, level LogLevel, tstamp time.Time)
}
