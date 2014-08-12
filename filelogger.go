package log4go

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	RollNone     RollFrequency = iota
	RollMinutely               // nice for testing
	RollHourly
	RollDaily
	RollWeekly
)

type RollFrequency uint8

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

func calculateNextRollTime(t time.Time, freq RollFrequency) time.Time {
	if freq == RollMinutely {
		t = t.Truncate(time.Minute)
		return t.Add(time.Minute)
	} else if freq == RollHourly {
		t = t.Truncate(time.Hour)
		return t.Add(time.Hour)
	} else {
		t = t.Truncate(time.Hour)
		day := t.Day()
		for ; t.Day() == day; t = t.Add(time.Hour) {
		}
		if freq == RollDaily {
			return t
		} else {
			for ; t.Day() != day; t = t.Add(time.Hour) {
			}
			return t
		}
	}
}

func calculatePreviousRollTime(t time.Time, freq RollFrequency) time.Time {
	if freq == RollMinutely {
		return t.Add(-time.Minute)
	} else if freq == RollHourly {
		return t.Add(-time.Hour)
	} else if freq == RollDaily {
		t = t.Add(-time.Hour * 12)
		for day := t.Day(); t.Day() == day; t = t.Add(-time.Hour) {
		}
		return t.Add(time.Hour)
	} else {
		t = t.Add(-time.Hour * 156)
		for day := t.Day(); t.Day() == day; t = t.Add(-time.Hour) {
		}
		return t.Add(time.Hour)
	}
}

func durationForRollFrequency(freq RollFrequency) time.Duration {
	var d time.Duration
	switch freq {
	case RollMinutely:
		d, _ = time.ParseDuration("1m")
	case RollHourly:
		d, _ = time.ParseDuration("1h")
	case RollDaily:
		d, _ = time.ParseDuration("24h")
	case RollWeekly:
		d, _ = time.ParseDuration("168h")
	}
	return d
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
