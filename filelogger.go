package log4go

import (
  "os"
  "path/filepath"
  "sync"
  "time"
)

const (
  defaultDirectoryPerms = 0755
  defaultFilePerms      = 0644
)

type fileAppender struct {
  lock sync.Mutex
  f *os.File
}

var fileAppenderMap = make(map[string]*fileAppender)
var fileAppenderMapLock = sync.Mutex{}

func (a *fileAppender) Append(msg string, level LogLevel, tstamp time.Time) {
  a.lock.Lock()
  defer a.lock.Unlock()

  a.f.Write([]byte(msg))
}

func NewFileLogger(dir string, filename string, level LogLevel, timePrefix string) (_ Log4Go, err error) {
  expandedDir, err := filepath.Abs(filepath.Dir(dir))
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
    appender = &fileAppender{sync.Mutex{}, logfile}
    fileAppenderMap[fullFilename] = appender
  }

  result := stdLogger{
    appender,
    level,
    timePrefix,
  }
  return Log4Go(&result), nil
}
