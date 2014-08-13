package log5go

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type logBuilder struct {
	level      LogLevel
	appender   Appender
	timeFormat string
	errs       *compositeError
}

// Entry point for building a new logger. Start here. Takes the desired log level.
func Log(level LogLevel) LogBuilder {
	builder := logBuilder{
		level,
		nil,
		TF_GoStd,
		newCompositeError(),
	}
	return &builder
}

// Add a custom format to the logger
func (b *logBuilder) WithTimeFmt(format string) LogBuilder {
	b.timeFormat = format
	return b
}

// Select the console appender. You must select an appender only once.
// You must select an appender prior to configuring it.
func (b *logBuilder) ToStdout() LogBuilder {
	if b.appender != nil {
		b.errs.append(fmt.Errorf("appender cannot be set more than once"))
	}

	b.appender = &consoleAppender{false}
	return b
}

// Select the file appender. You must select an appender only once.
// You must select an appender prior to configuring it.
func (b *logBuilder) ToFile(directory string, filename string) LogBuilder {
	if b.appender != nil {
		b.errs.append(fmt.Errorf("appender cannot be set more than once"))
	}

	expandedDir, err := filepath.Abs(directory)
	if err != nil {
		b.errs.append(err)
		return b
	}

	fullFilename := filepath.Join(expandedDir, filename)

	fileAppenderMapLock.Lock()
	defer fileAppenderMapLock.Unlock()

	var appender *fileAppender = fileAppenderMap[fullFilename]
	if appender == nil {
		logfile, err := os.Create(fullFilename)
		if err != nil {
			b.errs.append(err)
			return b
		}
		appender = &fileAppender{sync.Mutex{}, logfile, time.Now(), RollNone, SaveAllOldLogs}
		fileAppenderMap[fullFilename] = appender
	}

	if !fileRollerRunning {
		go periodicFileRoller()
		fileRollerRunning = true
	}

	b.appender = appender
	return b
}

// Add file rotation configuration to the file appender. ToFile() must have been
// called already.
func (b *logBuilder) WithRotation(frequency rollFrequency, keepNLogs int) LogBuilder {
	if b.appender == nil {
		b.errs.append(fmt.Errorf("appender must be set first"))
		return b
	}

	a, isFileAppender := b.appender.(*fileAppender)
	if !isFileAppender {
		b.errs.append(fmt.Errorf("appender not set to file appender"))
		return b
	}

	a.nextRollTime = calculateNextRollTime(time.Now(), frequency)
	a.rollFrequency = frequency
	a.keepNLogs = keepNLogs

	return b
}

// Send WARN, ERROR, and FATAL messages to stderr. ToConsole() must have been
// called already.
func (b *logBuilder) WithStderr() LogBuilder {
	if b.appender == nil {
		b.errs.append(fmt.Errorf("appender must be set first"))
		return b
	}

	a, isConsoleAppender := b.appender.(*consoleAppender)
	if !isConsoleAppender {
		b.errs.append(fmt.Errorf("appender not set to console appender"))
		return b
	}

	a.stderrAware = true
	return b
}

// Build the logger you have been configuring. Returns the logger, or any errors
// that have been encountered during the build process.
func (b *logBuilder) Build() (_ Log5Go, _ error) {
	if b.appender == nil {
		b.errs.append(fmt.Errorf("cannot build without appender set"))
	}

	if b.errs.hasErrors() {
		return nil, b.errs
	}

	logger := logger{
		b.level,
		b.appender,
		b.timeFormat,
	}
	return &logger, nil
}

// Build and register the logger you have been configuring. Returns the logger, or any errors
// that have been encountered during the build/register process.
func (b *logBuilder) Register(key string) (_ Log5Go, _ error) {
	logger, err := b.Build()
	if err != nil {
		return nil, err
	}

	err = loggerRegistry.Put(key, logger)
	if err != nil {
		return nil, err
	}

	return logger, nil
}
