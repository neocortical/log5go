package log4go

// Log4Go is log4go's primary logging interface. All logging is performed using
// the methods defined here.
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

// LogBuilder is the interface for building loggers.
type LogBuilder interface {
	WithTimeFmt(format string) LogBuilder
	ToConsole() LogBuilder
	ToFile(directory string, filename string) LogBuilder
	WithRotation(frequency rollFrequency, keepNLogs int) LogBuilder
	WithStderr() LogBuilder
	// WithLayout(pattern string) LogBuilder // TODO
	Build() (Log4Go, error)
	Register(key string) (Log4Go, error)
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
func GetLog(name string) (_ Log4Go, err error) {
	key, err := generateLoggerKey(name)
	if err != nil {
		return nil, err
	}

	return loggerRegistry.Get(key)
}

// Standard timestamp formats. You can use any format from the time package or
// roll your own.
const (
	TF_GoStd = "2006/01/02 15:04:05" // Default
	TF_NCSA  = "02/Jan/2006:15:04:05 -0700"
)
