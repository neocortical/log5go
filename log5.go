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
	GetLogLevel() LogLevel
	SetLogLevel(level LogLevel)
	GoLogger
	LogBuilder
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
