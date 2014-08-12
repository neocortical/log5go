package log4go

type LogLevel uint16

// Standard log levels (lifted directly from log4j)
const (
  LogAll   LogLevel = 0
  LogTrace LogLevel = 100
  LogDebug LogLevel = 200
  LogInfo  LogLevel = 300
  LogWarn  LogLevel = 400
  LogError LogLevel = 500
  LogFatal LogLevel = 600
)

// maps log levels to prefix strings describing each. extensible
var levelMap = map[LogLevel]string{
  LogAll:   "LOG",
  LogTrace: "TRACE",
  LogDebug: "DEBUG",
  LogInfo:  "INFO",
  LogWarn:  "WARN",
  LogError: "ERROR",
  LogFatal: "FATAL",
}

// Replace a log level prefix string, or add one for a custom log level
func RegisterLogLevel(level LogLevel, prefix string) {
  levelMap[level] = prefix
}
