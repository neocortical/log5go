package log4go

// Standard timestamp formats
const (
	TF_NCSA  = "02/Jan/2006:15:04:05 -0700"
	TF_GoStd = "2006/01/02 15:04:05"
)

type RollFrequency uint8

// Log rotation frequencies. Daily rotates at midnight, weekly rotates on Sunday at midnight
const (
	RollNone     RollFrequency = iota
	RollMinutely               // nice for testing
	RollHourly
	RollDaily
	RollWeekly
)

const SaveAllOldLogs = -1

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

type LogBuilder interface {
	WithTimeFormat(format string) LogBuilder
	ToConsole() LogBuilder
	ToFile(directory string, filename string) LogBuilder
	WithFileRotation(frequency RollFrequency, keepNLogs int) LogBuilder
	WithStderrSupport() LogBuilder
// WithLayout(pattern string) LogBuilder // TODO
	Build() (Log4Go, error)
//	BuildAndRegister(key string) (Log4Go, error) // TODO
}
