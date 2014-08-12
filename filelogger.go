package log4go

import (
	"os"
	"path/filepath"
	"sync"
	"time"
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

// Create a new file logger
func NewFileLogger(dir string, filename string, level LogLevel, timePrefix string) (_ Log4Go, err error) {
	return NewRollingFileLogger(dir, filename, level, timePrefix, RollNone, 0)
}

// Create a new rolling file logger that rotates logs with freq frequency
func NewRollingFileLogger(dir string, filename string, level LogLevel, timePrefix string, freq RollFrequency, oldLogsToSave int) (_ Log4Go, err error) {
	expandedDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	fullFilename := filepath.Join(expandedDir, filename)

	fileAppenderMapLock.Lock()
	defer fileAppenderMapLock.Unlock()

	if !fileRollerRunning {
		go periodicFileRoller()
		fileRollerRunning = true
	}

	var appender *fileAppender = fileAppenderMap[fullFilename]
	if appender == nil {
		logfile, err := os.Create(fullFilename)
		if err != nil {
			return nil, err
		}
		appender = &fileAppender{sync.Mutex{}, logfile, calculateNextRollTime(time.Now(), freq), freq, oldLogsToSave}
		fileAppenderMap[fullFilename] = appender
	}

	result := stdLogger{
		appender,
		level,
		timePrefix,
	}
	return Log4Go(&result), nil
}
