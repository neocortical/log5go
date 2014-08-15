package log5go

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Inner type of all loggers
type logger struct {
	lock       sync.Mutex
	level      LogLevel
	appender   Appender
	timeFormat string
	prefix     string
	lines      int
	flag       int 				// needed to return by Flags()
	buf        []byte     // for accumulating text to write
}

var std = initStd()

func initStd() (_ *logger) {
	log, _ := Log(LogAll).ToStderr().WithTimeFmt(TF_GoStd).Build()
	l, _ := log.(*logger)
	return l
}


var errLowLevel = errors.New("level too low")

// Log a message at the given log level
func (l *logger) Log(level LogLevel, format string, a ...interface{}) {
	l.log(time.Now(), level, 2, fmt.Sprintf(format, a...))
}

func (l *logger) Trace(format string, a ...interface{}) {
	l.log(time.Now(), LogTrace, 2, fmt.Sprintf(format, a...))
}

func (l *logger) Debug(format string, a ...interface{}) {
	l.log(time.Now(), LogDebug, 2, fmt.Sprintf(format, a...))
}

func (l *logger) Info(format string, a ...interface{}) {
	l.log(time.Now(), LogInfo, 2, fmt.Sprintf(format, a...))
}

func (l *logger) Warn(format string, a ...interface{}) {
	l.log(time.Now(), LogWarn, 2, fmt.Sprintf(format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	l.log(time.Now(), LogError, 2, fmt.Sprintf(format, a...))
}

func (l *logger) Fatal(format string, a ...interface{}) {
	l.log(time.Now(), LogFatal, 2, fmt.Sprintf(format, a...))
}

func (l *logger) GetLogLevel() LogLevel {
	return l.level
}

func (l *logger) SetLogLevel(level LogLevel) {
	l.level = level
}

func (l *logger) log(t time.Time, level LogLevel, calldepth int, s string) error {
	now := time.Now() // get this early.
	var file string
	var line int

	l.lock.Lock()
	defer l.lock.Unlock()

	if level < l.level {
		return errLowLevel
	}

	if l.lines != 0 {
		// release lock while getting caller info - it's expensive.
		l.lock.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.lock.Lock()
	}

	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, level, file, line)
	l.buf = append(l.buf, s...) // TODO: Appender should take []byte
	if len(s) > 0 && s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	l.appender.Append(string(l.buf), level, now)
	return nil // TODO: Appender should return error
}

func (l *logger) formatHeader(buf *[]byte, t time.Time, level LogLevel, file string, line int) {
	if l.timeFormat != "" {
		*buf = append(*buf, t.Format(l.timeFormat)...)
		*buf = append(*buf, ' ')
	}

	levelString := GetLogLevelString(level)
	if levelString != "" {
		*buf = append(*buf, levelString...)
		*buf = append(*buf, ' ')
	}

	if l.prefix != "" {
		*buf = append(*buf, l.prefix...)
	}

	if l.lines&(Lshortfile|Llongfile) != 0 {
		if l.lines == Lshortfile {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}

		if len(*buf) > 0 && (*buf)[len(*buf) - 1] != ' ' {
			*buf = append(*buf, ' ')
		}
		*buf = append(*buf, '(')
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		*buf = append(*buf, strconv.FormatUint(uint64(line), 10)...)
		*buf = append(*buf, ')')
	}

	// if no header info has been added, we don't print the separator
	if len(*buf) > 0 {
		*buf = append(*buf, ": "...)
	}
}
