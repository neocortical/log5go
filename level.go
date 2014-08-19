package log5go

import "sync"

type LogLevel uint16

// Standard log levels. Map to intergers separated by 100 to allow for custom
// log levels to be intermingled with standard ones.
const (
	LogAll LogLevel = iota * 100	// Log all messages, regardless of level
	LogTrace											// TRACE log leve/threshold
	LogDebug											// DEBUG log leve/threshold
	LogInfo												// INFO log leve/threshold
	LogWarn												// WARN log leve/threshold
	LogError											// ERROR log leve/threshold
	LogFatal											// FATAL log leve/threshold
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

// Get the string value of a LogLevel if it is registered
func GetLogLevelString(level LogLevel) string {
	levelMapLock.RLock()
	defer levelMapLock.RUnlock()
	return levelMap[level]
}
