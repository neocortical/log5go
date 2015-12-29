package log5go

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type fileAppender struct {
	lock          sync.Mutex
	f             *os.File
	fname         string // full path to file
	lastOpenTime  time.Time
	nextRollTime  time.Time
	rollFrequency rollFrequency
	keepNLogs     int
}

var fileAppenderMap = make(map[string]*fileAppender)
var fileAppenderMapLock = sync.Mutex{}
var fileRollerRunning = false

var fileWatcherPeriod time.Duration = time.Second
var fileReopenRefractoryPeriod time.Duration = time.Second

func (a *fileAppender) Append(msg *[]byte, level LogLevel, tstamp time.Time) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	TerminateMessageWithNewline(msg)

	if a.shouldRoll(tstamp) {
		a.doRoll()
	}

	// sanity check here, in case file couldn't be reopened after rolling
	if a.f == nil {
		return fmt.Errorf("file couldn't be opened")
	}
	_, err := a.f.Write(*msg)
	return err
}

// Determine whether we should roll the log file. Must be in lock already.
func (a *fileAppender) shouldRoll(tstamp time.Time) bool {
	if a.rollFrequency == RollNone {
		return false
	} else {
		return !tstamp.Before(a.nextRollTime)
	}
}

// Actually roll the log file. Must be in lock already.
func (a *fileAppender) doRoll() {
	absoluteFilename := a.fname
	dir, filename := filepath.Split(absoluteFilename)
	a.f.Close()

	archiveTime := calculatePreviousRollTime(a.nextRollTime, a.rollFrequency)
	archiveFilename := generateArchiveFilename(filename, archiveTime, a.rollFrequency)
	a.nextRollTime = calculateNextRollTime(a.nextRollTime, a.rollFrequency)

	archiveAbsFilename := filepath.Join(dir, archiveFilename)
	os.Rename(absoluteFilename, archiveAbsFilename)
	a.f, _ = os.Create(absoluteFilename)

	// if we are saving N archived logs, try to delete N+1
	if a.keepNLogs > -1 {
		for i := 0; i < a.keepNLogs; i++ {
			archiveTime = calculatePreviousRollTime(archiveTime, a.rollFrequency)
		}
		deleteFilename := filepath.Join(dir, generateArchiveFilename(filename, archiveTime, a.rollFrequency))
		os.Remove(deleteFilename)
	}
}

func generateArchiveFilename(fname string, rollTime time.Time, freq rollFrequency) string {
	var timeFormat string
	if freq == RollMinutely || freq == RollHourly {
		timeFormat = "2006-01-02-15-04-MST"
	} else {
		timeFormat = "2006-01-02"
	}

	return fmt.Sprintf("%s.%s", fname, rollTime.Format(timeFormat))
}

// validFile method returns true iff file is open and valid. Caller should protect with lock.
func (a *fileAppender) validFile() bool {
	_, err := os.Stat(a.f.Name())
	return err == nil
}

// shouldAttemptFileReopen method returns true iff the appender should attempt to reopen file
// Caller should protect by locking appender.
func (a *fileAppender) shouldAttemptFileReopen() bool {
	return time.Now().Sub(a.lastOpenTime) > fileReopenRefractoryPeriod
}

func (a *fileAppender) reopenFile() {
	a.lastOpenTime = time.Now()
	newFile, err := os.OpenFile(a.fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err == nil {
		a.f = newFile
	} else {
		a.f.Close() // ignore errors
	}
}

// run in goroutine to periodically check logs for rollability
func periodicFileWatcher() {
	ticker := time.NewTicker(fileWatcherPeriod)
	for {
		tick := <-ticker.C

		watchFiles(tick)
	}
}

func watchFiles(tick time.Time) {
	fileAppenderMapLock.Lock()
	for _, a := range fileAppenderMap {
		a.lock.Lock()
		if !a.validFile() {
			if a.shouldAttemptFileReopen() {
				a.reopenFile()
			} else {
				a.f.Close() // ignore errors
			}
		}
		if a.shouldRoll(tick) {
			a.doRoll()
		}
		a.lock.Unlock()
	}
	fileAppenderMapLock.Unlock()
}
