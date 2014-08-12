package log4go

import (
  "path/filepath"
  "fmt"
  "os"
  "sync"
  "time"
)

type logBuilder struct {
  level LogLevel
  appender appender
  timeFormat string
  errs *compositeError
}

func NewLog(level LogLevel) LogBuilder {
  builder := logBuilder{
    level,
    nil,
    TF_GoStd,
    newCompositeError(),
  }
  return &builder
}

func (b *logBuilder) WithTimeFormat(format string) LogBuilder {
  b.timeFormat = format
  return b
}

func (b *logBuilder) ToConsole() LogBuilder {
  if b.appender != nil {
    b.errs.append(fmt.Errorf("appender cannot be set more than once"))
  }

  b.appender = &consoleAppender{false}
  return b
}

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
    appender = &fileAppender{sync.Mutex{}, logfile, time.Now(), RollNone, -1}
    fileAppenderMap[fullFilename] = appender
  }

  if !fileRollerRunning {
    go periodicFileRoller()
    fileRollerRunning = true
  }

  b.appender = appender
  return b
}

func (b *logBuilder) WithFileRotation(frequency RollFrequency, keepNLogs int) LogBuilder {
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

func (b *logBuilder) WithStderrSupport() LogBuilder {
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

func (b *logBuilder) Build() (_ Log4Go, _ error) {
  if b.appender == nil {
    b.errs.append(fmt.Errorf("cannot build without appender set"))
  }

  if b.errs.hasErrors() {
    return nil, b.errs
  }

  logger := stdLogger{
    b.level,
    b.appender,
    b.timeFormat,
  }
  return &logger, nil
}
