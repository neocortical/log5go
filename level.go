package log5go

import "sync"

type LogLevel uint16

// Standard log levels. Map to intergers separated by 100 to allow for custom
// log levels to be intermingled with standard ones.
const (
	LogAll      LogLevel = 0   // Log all messages, regardless of level
	LogTrace    LogLevel = 100 // TRACE log leve/threshold
	LogDebug    LogLevel = 200 // DEBUG log leve/threshold
	LogInfo     LogLevel = 300 // INFO log leve/threshold
	LogNotice   LogLevel = 350 // NOTICE log leve/threshold
	LogWarn     LogLevel = 400 // WARN log leve/threshold
	LogError    LogLevel = 500 // ERROR log leve/threshold
	LogCritical LogLevel = 530 // CRITICAL log leve/threshold
	LogAlert    LogLevel = 560 // ALERT log leve/threshold
	LogFatal    LogLevel = 600 // FATAL/EMERG log leve/threshold
)

// maps log levels to prefix strings describing each. extensible
var levelMap = map[LogLevel]string{
	LogAll:      "ALL",
	LogTrace:    "TRACE",
	LogDebug:    "DEBUG",
	LogInfo:     "INFO",
	LogNotice:   "NOTICE",
	LogWarn:     "WARN",
	LogError:    "ERROR",
	LogCritical: "CRIT",
	LogAlert:    "ALERT",
	LogFatal:    "FATAL", // This level is also Syslog/EMERG
}

// Protects levelMap
var levelMapLock = new(sync.RWMutex)

// RegisterLogLevel replaces a log level prefix string, or adds one for a custom log level
func RegisterLogLevel(level LogLevel, prefix string) {
	levelMapLock.Lock()
	levelMap[level] = prefix
	levelMapLock.Unlock()
}

// DeregisterLogLevel deletes a log level string, so it will not be used in log messages at that level
func DeregisterLogLevel(level LogLevel) {
	levelMapLock.Lock()
	delete(levelMap, level)
	levelMapLock.Unlock()
}

// GetLogLevelString gets the string value of a LogLevel if it is registered
func GetLogLevelString(level LogLevel) string {
	levelMapLock.RLock()
	defer levelMapLock.RUnlock()
	return levelMap[level]
}

// LogLevelForString returns a LogLevel whose string matches the specified string.
// Duplicate level strings will produce indeterminate results. Returns LogAll if string not found.
func LogLevelForString(val string) LogLevel {
	levelMapLock.Lock()
	defer levelMapLock.Unlock()
	for k, v := range levelMap {
		if v == val {
			return k
		}
	}
	return LogAll
}
