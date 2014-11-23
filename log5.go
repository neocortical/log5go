package log5go

import (
	"io"
)

// Log5Go is log5go's primary logging interface. All logging is performed using
// the methods defined here.
type Log5Go interface {
	// Log logs a message at a custom log level (or explicitly at a standard log level)
	Log(level LogLevel, format string, a ...interface{})

	// Trace logs a message at the TRACE log level
	Trace(format string, a ...interface{})

	// Debug logs a message at the TRACE log level
	Debug(format string, a ...interface{})

	// Info logs a message at the INFO log level
	Info(format string, a ...interface{})

	// Warn logs a message at the WARN log level
	Warn(format string, a ...interface{})

	// Error logs a message at the ERROR log level
	Error(format string, a ...interface{})

	// Fatal logs a message at the FATAL log level. Note: Fatal() DOES NOT call os.Exit or panic.
	Fatal(format string, a ...interface{})

	// LogLevel returns the threshold that log messages must meet to be logged
	LogLevel() LogLevel

	// SetLogLevel sets the threshold that log messages must meet to be logged
	SetLogLevel(level LogLevel)

	// Go pkg/log compatibility. See the GoLogger interface for details
	GoLogger

	// LogBuilder contains methods for creating new logs using a builder pattern. See the LogBuilder interface for details.
	LogBuilder

	// Log5GoData contains methods for appending structured data to log messages. See interface description for details.
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

	// WithTimeFmt sets the time format that the logger will use. Use "" for no timestamp.
	WithTimeFmt(format string) Log5Go

	// ToStdout creates a logger that logs all messages to os.Stdout.
	ToStdout() Log5Go

	// ToStderr creates a logger that logs all messages to os.Stderr. This is the default.
	ToStderr() Log5Go

	// ToWriter creates a logger the writes to the specified Writer. Concurrency-safe writes.
	ToWriter(out io.Writer) Log5Go

	// ToFile creates a logger that appends to the specified file. Currently, file is deleted if it exists.
	ToFile(directory string, filename string) Log5Go

	// ToAppender creates a logger that appends to a user-supplied appender.
	ToAppender(appender Appender) Log5Go

	// WithRotation sets file rotation information for a logger set to append to a file with ToFile()
	WithRotation(frequency rollFrequency, keepNLogs int) Log5Go

	// WithStderr writes all log messages at WARN or above to os.Stderr
	WithStderr() Log5Go

	// WithPrefix sets a custom prefix that will appear in all logged messages
	WithPrefix(prefix string) Log5Go

	// WithLine adds full caller information (/full/path/to/file.go:linenum) to logged messages
	WithLine() Log5Go

	// WithLn adds partial caller info (file.go:linenum) to logged messages
	WithLn() Log5Go

	// WithFmt sets a custom string format for log messages. See StringFormatter for details
	WithFmt(format string) Log5Go

	// Json causes all log messages to JSON-formatted.
	Json() Log5Go

	// Register registers a logger in the log5go registry, allowing it to be retrieved from anywhere in your program
	Register(key string) (Log5Go, error)

	SetForApplication() Log5Go
}

type rollFrequency uint8

// Log rotation frequencies. Daily rotates at midnight, weekly rotates on Sunday at midnight
const (
	RollNone     rollFrequency = iota // Don't do file rotation
	RollMinutely                      // Rotate files once per minute
	RollHourly                        // Rotate files once per hour
	RollDaily                         // Rotate files once per day
	RollWeekly                        // Rotate files once per week
)

// SaveAllLogs used as an argument to WithFileRotation(, keepNLogs). Disables deleting of old log files.
const SaveAllLogs = -1

// Gets a log by looking it up by name in the internal registry.
func GetLog(key string) (_ Log5Go, err error) {
	return loggerRegistry.Get(key)
}

func GetAppLogger() Log5Go {
	return applicationLogger
}

// Standard timestamp formats. You can use any format from the time package or
// roll your own.
const (
	TF_GoStd = "2006/01/02 15:04:05"        // Default
	TF_NCSA  = "02/Jan/2006:15:04:05 -0700" // NCSA standard time format
)
