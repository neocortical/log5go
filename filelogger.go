package log4go

import (
	"fmt"
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

type fileAppender struct {
	lock          sync.Mutex
	f             *os.File
	nextRollTime  time.Time
	rollFrequency RollFrequency
	rollsToSave   uint
}

var fileAppenderMap = make(map[string]*fileAppender)
var fileAppenderMapLock = sync.Mutex{}

func (a *fileAppender) Append(msg string, level LogLevel, tstamp time.Time) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.shouldRoll(tstamp) {
		a.doRoll()
	}

	// sanity check here, in case file couldn't be reopened after rolling
	if a.f != nil {
		a.f.Write([]byte(msg))
	}
}

func (a *fileAppender) shouldRoll(tstamp time.Time) bool {
	if a.rollFrequency == RollNone {
		return false
	} else {
		return !tstamp.Before(a.nextRollTime)
	}
}

func (a *fileAppender) doRoll() {
	absoluteFilename := a.f.Name()
	dir, filename := filepath.Split(absoluteFilename)
	a.f.Close()

	lastTime := calculatePreviousRollTime(a.nextRollTime, a.rollFrequency)
	a.nextRollTime = calculateNextRollTime(a.nextRollTime, a.rollFrequency)

	var timeFormat string
	if a.rollFrequency == RollMinutely || a.rollFrequency == RollHourly {
		timeFormat = "2006-01-02T15-04"
	} else {
		timeFormat = "2006-01-02"
	}

	archiveFilename := filepath.Join(dir, fmt.Sprintf("%s.%s", filename, lastTime.Format(timeFormat)))
	os.Rename(absoluteFilename, archiveFilename)
	a.f, _ = os.Create(absoluteFilename)
}

// Create a new file logger
func NewFileLogger(dir string, filename string, level LogLevel, timePrefix string) (_ Log4Go, err error) {
	return NewRollingFileLogger(dir, filename, level, timePrefix, RollNone, 0)
}

// Create a new rolling file logger that rotates logs with freq frequency
func NewRollingFileLogger(dir string, filename string, level LogLevel, timePrefix string, freq RollFrequency, oldLogsToSave uint) (_ Log4Go, err error) {
	expandedDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	fullFilename := filepath.Join(expandedDir, filename)

	fileAppenderMapLock.Lock()
	defer fileAppenderMapLock.Unlock()

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
