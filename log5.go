package log5go

import (
	"io"
)

// Log5Go is log5go's primary logging interface. All logging is performed using
// the methods defined here.
type Log5Go interface {
	Log(level LogLevel, format string, a ...interface{})
	Trace(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
	LogLevel() LogLevel
	SetLogLevel(level LogLevel)
	GoLogger
	LogBuilder
	Log5GoData
}

// Log5GoData interface allows developers to add custom structured data to log
// messages with WithData(d Data). WithData() is intended to be called immediately
// before calling a log method. Developers should not attempt to modify the
// configuration of a logger after calling WithData().
//
// Data is an alias for map[string]interface{}. Callers should only insert builtin
// types into the data map; any non-builtin types are scrubbed for compaibility with
// different formatters.
type Log5GoData interface {
  WithData(d Data) Log5Go
}

// LogBuilder is the interface for building loggers.
type LogBuilder interface {
	WithTimeFmt(format string) Log5Go
	ToStdout() Log5Go
	ToStderr() Log5Go
	ToWriter(out io.Writer) Log5Go
	ToFile(directory string, filename string) Log5Go
	ToAppender(appender Appender) Log5Go
	WithRotation(frequency rollFrequency, keepNLogs int) Log5Go
	WithStderr() Log5Go
	WithPrefix(prefix string) Log5Go
	WithLine() Log5Go
	WithLn() Log5Go
	// With a custom string format
	WithFmt(format string) Log5Go
	Json() Log5Go
	Register(key string) (Log5Go, error)
}

type rollFrequency uint8

// Log rotation frequencies. Daily rotates at midnight, weekly rotates on Sunday at midnight
const (
	RollNone     rollFrequency = iota
	RollMinutely               // nice for testing
	RollHourly
	RollDaily
	RollWeekly
)

// SaveAllOldLogs used as an argument to WithFileRotation(, keepNLogs)
const SaveAllOldLogs = -1

// Gets a log by looking it up by name in the internal registry.
func GetLog(key string) (_ Log5Go, err error) {
	return loggerRegistry.Get(key)
}

// Standard timestamp formats. You can use any format from the time package or
// roll your own.
const (
	TF_GoStd = "2006/01/02 15:04:05" // Default
	TF_NCSA  = "02/Jan/2006:15:04:05 -0700"
)
