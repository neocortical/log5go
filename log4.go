package log4go

const (
	TF_NCSA  = "02/Jan/2006:15:04:05 -0700"
	TF_GoStd = "2006/01/02 15:04:05"
)

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
