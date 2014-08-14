package log5go

import (
	"time"
)

// Appender defines the interface responsible for writing a processed message to
// a destination. Currently, destinations are console and file. Additional
// appender types are planned.
type Appender interface {
	Append(msg string, level LogLevel, tstamp time.Time)
}

// TODO: Appender should take []byte
// TODO: Appender should return error (to satisfy compat Output)
