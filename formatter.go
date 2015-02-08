package log5go

import "time"

// interface Formatter formats a log message into *out for passing to an Appender
type Formatter interface {
	Format(tstamp time.Time, level LogLevel, prefix, caller string, line uint, msg string, data Data, out *[]byte)
	SetTimeFormat(timeFormat string)
	SetLines(lines bool)
}

// Some constant string formats for convenience, also used internally
const (
	FMT_Default            = "%t %l : %m"
	FMT_DefaultPrefix      = "%t %l %p: %m"
	FMT_DefaultLines       = "%t %l (%c:%n): %m"
	FMT_DefaultPrefixLines = "%t %l %p (%c:%n): %m"
	FMT_NoTime             = "%l : %m"
	FMT_NoTimePrefix       = "%l %p: %m"
	FMT_NoTimeLines        = "%l (%c:%n): %m"
	FMT_NoTimePrefixLines  = "%l %p (%c:%n): %m"
)
