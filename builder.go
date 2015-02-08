package log5go

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

// Logger is the entry point for building a new logger. Takes the desired log level threshold and returns a stderr logger.
func Logger(level LogLevel) Log5Go {
	logger := logger{
		level:      level,
		formatter:  NewStringFormatter(FMT_Default),
		appender:   &writerAppender{dest: os.Stderr, errDest: nil},
		timeFormat: TF_GoStd,
		prefix:     "",
		lines:      0,
	}
	return &logger
}

func (l *logger) Clone() Log5Go {
	return &logger{
		level:      l.level,
		formatter:  l.formatter,
		appender:   l.appender,
		timeFormat: l.timeFormat,
		prefix:     l.prefix,
		lines:      l.lines,
	}
}

// Add a custom format to the logger
func (l *logger) WithTimeFmt(format string) Log5Go {
	l.timeFormat = format
	l.formatter.SetTimeFormat(format)
	return l
}

// Select the console appender set to stdout. You must select an appender only once.
// You must select an appender prior to configuring it.
func (l *logger) ToStdout() Log5Go {
	l.appender = &writerAppender{dest: os.Stdout, errDest: nil}
	return l
}

// Select the console appender set to stderr. You must select an appender only once.
// You must select an appender prior to configuring it.
func (l *logger) ToStderr() Log5Go {
	l.appender = &writerAppender{dest: os.Stderr, errDest: nil}
	return l
}

// Select the console appender with a custom destination.
// You must select an appender only once.
// You must select an appender prior to configuring it.
func (l *logger) ToWriter(out io.Writer) Log5Go {
	l.appender = &writerAppender{dest: out, errDest: nil}
	return l
}

func (l *logger) WithPrefix(prefix string) Log5Go {
	l.prefix = prefix
	l.updateFormatterIfNecessary()
	return l
}

func (l *logger) WithLongLines() Log5Go {
	l.lines = LogLinesLong
	l.updateFormatterIfNecessary()
	return l
}

func (l *logger) WithShortLines() Log5Go {
	l.lines = LogLinesShort
	l.updateFormatterIfNecessary()
	return l
}

// Select the file appender. You must select an appender only once.
// You must select an appender prior to configuring it.
func (l *logger) ToFile(directory string, filename string) Log5Go {
	expandedDir, err := filepath.Abs(directory)
	if err != nil {
		// would be nice to do *something* on error, but not sure what
		return l
	}

	fullFilename := filepath.Join(expandedDir, filename)

	fileAppenderMapLock.Lock()
	var appender = fileAppenderMap[fullFilename]
	if appender == nil {
		logfile, err := os.OpenFile(fullFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fileAppenderMapLock.Unlock()
			return l
		}
		appender = &fileAppender{
			f:             logfile,
			fname:         fullFilename,
			lastOpenTime:  time.Now(),
			nextRollTime:  time.Now(),
			rollFrequency: RollNone,
			keepNLogs:     SaveAllLogs,
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
	l.appender = appender
	return l
}

// ToLocalSyslog sets a syslog formatter and attempts to set a syslog appender connected
// to the local syslogd daemon. If this fails, stderr is used instead and an error message
// is immediately logged.
func (l *logger) ToLocalSyslog(facility SyslogPriority, tag string) Log5Go {

	if facility < SyslogKernel || facility > SyslogLocal7 {
		l.appender = &writerAppender{dest: os.Stderr, errDest: os.Stderr}
		l.Error("INVALID SYSLOG FACILITY: %d", facility)

		return l
	}

	var conn net.Conn
	var err error
	for _, transport := range socketTypes {
		for _, socket := range socketLocations {
			conn, err = net.DialTimeout(transport, socket, time.Second*10)
			if err != nil {
				continue
			} else {
				l.appender = &syslogAppender{conn: conn, facility: facility, tag: tag}
				l.formatter = newSyslogFormatter(l.lines != 0)
				return l
			}
		}
	}

	l.appender = &writerAppender{dest: os.Stderr, errDest: os.Stderr}
	l.Error("UNABLE TO CONNECT TO LOCAL SYSLOG PROCESS: %v", err)

	return l
}

// ToRemoteSyslog sets a syslog formatter and attempts to set a syslog appender connected
// to the remote syslogd daemon. If this fails, stderr is used instead and an error message
// is immediately logged.
func (l *logger) ToRemoteSyslog(facility SyslogPriority, tag string, transport string, addr string) Log5Go {

	if facility < SyslogKernel || facility > SyslogLocal7 {
		l.appender = &writerAppender{dest: os.Stderr, errDest: os.Stderr}
		l.Error("INVALID SYSLOG FACILITY: %d", facility)

		return l
	}

	var conn net.Conn
	var err error
	fmt.Printf("dialing %s (%s).\n", addr, transport)

	conn, err = net.DialTimeout(transport, addr, time.Second*10)
	fmt.Printf("finished dialing. error: %v\n", err)
	if err == nil {
		l.appender = &syslogAppender{conn: conn, facility: facility, tag: tag}
		l.formatter = newSyslogFormatter(l.lines != 0)
		return l
	}

	l.appender = &writerAppender{dest: os.Stderr, errDest: os.Stderr}
	l.Error("UNABLE TO CONNECT TO REMOTE SYSLOG PROCESS: %v", err)

	return l
}

// Add file rotation configuration to the file appender. ToFile() must have been
// called already.
func (l *logger) WithRotation(frequency rollFrequency, keepNLogs int) Log5Go {
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
	a, iswriterAppender := l.appender.(*writerAppender)
	if !iswriterAppender {
		return l
	}

	a.errDest = os.Stderr
	return l
}

func (l *logger) WithFmt(format string) Log5Go {
	stringFormatter := NewStringFormatter(format)
	stringFormatter.explicitFormat = true
	l.formatter = stringFormatter
	l.formatter.SetTimeFormat(l.timeFormat)
	l.formatter.SetLines(l.lines != 0)
	return l
}

// Build and register the logger you have been configuring. Returns the logger, or any errors
// that have been encountered during the build/register process.
func (l *logger) Register(key string) Log5Go {
	loggerRegistry.Put(key, l)
	return l
}

// getDefaultFormat method inspects the logger and applies the appropriate default
// format for the current config. logger should be locked by the caller so that
// config remains unchained when the data is rendered for the returned format.
func (l *logger) updateFormatterIfNecessary() {
	stringFormatter, ok := l.formatter.(*StringFormatter)
	if ok && !stringFormatter.explicitFormat {
		pattern := getFormatForSettings(l.prefix, l.timeFormat, l.lines != 0)
		stringFormatter.parts = decodePattern(pattern)
	}

	l.formatter.SetTimeFormat(l.timeFormat)
	l.formatter.SetLines(l.lines != 0)
}

func getFormatForSettings(prefix, timeFormat string, lines bool) string {
	if timeFormat == "" {
		if lines {
			if prefix != "" {
				return FMT_DefaultPrefixLines
			} else {
				return FMT_NoTimeLines
			}
		} else if prefix != "" {
			return FMT_NoTimePrefix
		} else {
			return FMT_NoTime
		}
	} else {
		if lines {
			if prefix != "" {
				return FMT_DefaultPrefixLines
			} else {
				return FMT_DefaultLines
			}
		} else if prefix != "" {
			return FMT_DefaultPrefix
		} else {
			return FMT_Default
		}
	}
}
