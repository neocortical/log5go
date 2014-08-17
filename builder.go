package log5go

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

// Entry point for building a new logger. Start here. Takes the desired log level.
func Logger(level LogLevel) Log5Go {
	logger := logger{
		level: level,
		formatter: nil,
		appender: &writerAppender{dest: os.Stderr, errDest: nil},
		timeFormat: TF_GoStd,
		prefix: "",
		lines: 0,
	}
	return &logger
}

// Add a custom format to the logger
func (l *logger) WithTimeFmt(format string) Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.timeFormat = format
	return l
}

// Select the console appender set to stdout. You must select an appender only once.
// You must select an appender prior to configuring it.
func (l *logger) ToStdout() Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.appender = &writerAppender{dest: os.Stdout, errDest: nil}
	return l
}

// Select the console appender set to stderr. You must select an appender only once.
// You must select an appender prior to configuring it.
func (l *logger) ToStderr() Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.appender = &writerAppender{dest: os.Stderr, errDest: nil}
	return l
}

// Select the console appender with a custom destination.
// You must select an appender only once.
// You must select an appender prior to configuring it.
func (l *logger) ToWriter(out io.Writer) Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.appender = &writerAppender{dest: out, errDest: nil}
	return l
}

func (l *logger) WithPrefix(prefix string) Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.prefix = prefix
	return l
}

func (l *logger) WithLine() Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lines = Llongfile
	return l
}

func (l *logger) WithLn()  Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lines = Lshortfile
	return l
}

// Select the file appender. You must select an appender only once.
// You must select an appender prior to configuring it.
func (l *logger) ToFile(directory string, filename string) Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	expandedDir, err := filepath.Abs(directory)
	if err != nil {
		// would be nice to do *something* on error, but not sure what
		return l
	}

	fullFilename := filepath.Join(expandedDir, filename)

	fileAppenderMapLock.Lock()
	var appender *fileAppender = fileAppenderMap[fullFilename]
	if appender == nil {
		logfile, err := os.Create(fullFilename)
		if err != nil {
			fileAppenderMapLock.Unlock()
			return l
		}
		appender = &fileAppender{
			f: logfile,
			lastOpenTime: time.Now(),
			nextRollTime: time.Now(),
			rollFrequency: RollNone,
			keepNLogs: SaveAllLogs,
		}
		fileAppenderMap[fullFilename] = appender
	}

	if !fileRollerRunning {
		go periodicFileWatcher()
		fileRollerRunning = true
	}
	fileAppenderMapLock.Unlock()

	l.appender = appender
	return l
}

// ToAppender sets a custom (i.e third-party) appender as the destination for this logger.
// No other appender setting methods must be called before or after.
func (l *logger) ToAppender(appender Appender) Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.appender = appender
	return l
}

// Add file rotation configuration to the file appender. ToFile() must have been
// called already.
func (l *logger) WithRotation(frequency rollFrequency, keepNLogs int) Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()

	a, isFileAppender := l.appender.(*fileAppender)
	if !isFileAppender {
		return l
	}

	a.nextRollTime = calculateNextRollTime(time.Now(), frequency)
	a.rollFrequency = frequency
	a.keepNLogs = keepNLogs

	return l
}

// Send WARN, ERROR, and FATAL messages to stderr. ToConsole() must have been
// called already.
func (l *logger) WithStderr() Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()

	a, iswriterAppender := l.appender.(*writerAppender)
	if !iswriterAppender {
		return l
	}

	a.errDest = os.Stderr
	return l
}

func (l *logger) WithFmt(format string) Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.formatter = NewStringFormatter(format)
	return l
}

// Build and register the logger you have been configuring. Returns the logger, or any errors
// that have been encountered during the build/register process.
func (l *logger) Register(key string) (_ Log5Go, _ error) {
	err := loggerRegistry.Put(key, l)
	return l, err
}
