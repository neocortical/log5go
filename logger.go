package log5go

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
)

// Inner type of all loggers
type logger struct {
	sync.RWMutex
	level      LogLevel
	formatter  Formatter
	appender   Appender
	timeFormat string
	prefix     string
	lines      LogLines
	buf        []byte // buffer for holding formatted log messages
}

type LogLines int

const (
	LogLinesNone  LogLines = 0
	LogLinesShort LogLines = 1
	LogLinesLong  LogLines = 2
)

var std = initStd()

func initStd() (_ *logger) {
	log := Logger(LogAll).ToStderr().WithTimeFmt(TF_GoStd)
	l, _ := log.(*logger)
	return l
}

var errLowLevel = errors.New("level too low")

// Log a message at the given log level
func (l *logger) Log(level LogLevel, format string, a ...interface{}) {
	l.log(time.Now(), level, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) Trace(format string, a ...interface{}) {
	l.log(time.Now(), LogTrace, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) Debug(format string, a ...interface{}) {
	l.log(time.Now(), LogDebug, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) Info(format string, a ...interface{}) {
	l.log(time.Now(), LogInfo, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) Notice(format string, a ...interface{}) {
	l.log(time.Now(), LogNotice, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) Warn(format string, a ...interface{}) {
	l.log(time.Now(), LogWarn, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) Error(format string, a ...interface{}) {
	l.log(time.Now(), LogError, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) Critical(format string, a ...interface{}) {
	l.log(time.Now(), LogCritical, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) Alert(format string, a ...interface{}) {
	l.log(time.Now(), LogAlert, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) Fatal(format string, a ...interface{}) {
	l.log(time.Now(), LogFatal, 2, fmt.Sprintf(format, a...), nil)
}

func (l *logger) LogLevel() LogLevel {
	return l.level
}

func (l *logger) SetLogLevel(level LogLevel) {
	l.level = level
}

func (l *logger) WithData(d Data) Log5Go {
	return &boundLogger{l: l, data: d}
}

func (l *logger) Json() Log5Go {
	l.formatter = &jsonFormatter{}
	l.formatter.SetTimeFormat(l.timeFormat)
	l.formatter.SetLines(l.lines != 0)
	return l
}

// log method is the actual logging implementation. It takes all data about a logging
// event, prepares it, applies the appropriate formatter, and sends the data to the
// configured log appender.
func (l *logger) log(t time.Time, level LogLevel, calldepth int, msg string, data Data) error {
	now := time.Now() // get this early.
	var file string
	var line int

	if level < l.level {
		return errLowLevel
	}

	if l.lines != LogLinesNone {
		// release lock while getting caller info - it's expensive.
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}

		if l.lines == LogLinesShort {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
	}

	data = scrubData(data)

	// lock buffer
	l.Lock()
	defer l.Unlock()

	l.buf = l.buf[:0]
	l.formatter.Format(now, level, l.prefix, file, uint(line), msg, data, &l.buf)

	return l.appender.Append(&l.buf, level, now)
}

// scrubData scrubs map of any non-basic elements
func scrubData(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		if value == nil {
			continue // null values OK
		}
		switch reflect.TypeOf(value).Kind() {
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
			// let it through
		default:
			delete(data, key)
		}
	}
	return data
}
