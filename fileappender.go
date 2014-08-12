package log4go

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
  nextRollTime  time.Time
  rollFrequency RollFrequency
  rollsToSave   uint
}

var fileAppenderMap = make(map[string]*fileAppender)
var fileAppenderMapLock = sync.Mutex{}
var fileRollerRunning = false

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

func periodicFileRoller() {

	ticker := time.NewTicker(time.Second*15)
	for {
		tick := <-ticker.C

    fileAppenderMapLock.Lock()
    for _, a := range fileAppenderMap {
      a.lock.Lock()
      if a.shouldRoll(tick) {
        a.doRoll()
      }
      a.lock.Unlock()
    }
    fileAppenderMapLock.Unlock()
  }

}
