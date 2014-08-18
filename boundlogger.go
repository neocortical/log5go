package log5go

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Data represents user-added key/value pairs to a log message. For string output,
// these values are added to the end of the end of the user message with the form
// key=value. For JSON output, data is added as a JSON object, like
// data:{ "key1":10, "key2":"foo" }. Developers should avoid using complex types
// as data values, as it could break log output if the value cannot be properly
// formatted.
type Data map[string]interface{}

// boundLogger binds a logger to user-supplied data. boundLogger implements the Log5Go
// interface so a logging method can be called on it to log the data. The boundLogger
// object can be reused. Calling LogBuilder methods on a boundLogger object result in
// a NOOP to keep developers from doing silly things.
type boundLogger struct {
	l *logger
	lock sync.RWMutex // protects data
	data Data
}

//-- Log5GoData interface ------------

func (l *boundLogger) WithData(d Data) Log5Go {
	l.lock.Lock()
	defer l.lock.Unlock()
	for key, value := range d {
		l.data[key] = value
	}
	return l
}

//-- Log5Go interface ------------

func (l *boundLogger) Log(level LogLevel, format string, a ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), level, 2, fmt.Sprintf(format, a...), l.data)
}

func (l *boundLogger) Trace(format string, a ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogTrace, 2, fmt.Sprintf(format, a...), l.data)
}

func (l *boundLogger) Debug(format string, a ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogDebug, 2, fmt.Sprintf(format, a...), l.data)
}

func (l *boundLogger) Info(format string, a ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogInfo, 2, fmt.Sprintf(format, a...), l.data)
}

func (l *boundLogger) Warn(format string, a ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogWarn, 2, fmt.Sprintf(format, a...), l.data)
}

func (l *boundLogger) Error(format string, a ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogError, 2, fmt.Sprintf(format, a...), l.data)
}

func (l *boundLogger) Fatal(format string, a ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogFatal, 2, fmt.Sprintf(format, a...), l.data)
}

func (l *boundLogger) LogLevel() LogLevel {
	return l.l.LogLevel()
}

func (l *boundLogger) SetLogLevel(level LogLevel) {
	// NOOP
}

//-- GoLogger interface ------------------

func (l *boundLogger) Output(calldepth int, s string) error {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.l.log(time.Now(), LogInfo, calldepth + 1, s, l.data)
}

func (l *boundLogger) SetOutput(out io.Writer) {
	// NOOP
}

func (l *boundLogger) Flags() int {
	return l.l.Flags()
}

func (l *boundLogger) SetFlags(flag int) {
	// NOOP
}

func (l *boundLogger) Prefix() string {
	return l.l.Prefix()
}

func (l *boundLogger) SetPrefix(prefix string) {
	// NOOP
}

func (l *boundLogger) Printf(format string, v ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogInfo, 2, fmt.Sprintf(format, v...), l.data)
}

func (l *boundLogger) Print(v ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogInfo, 2, fmt.Sprint(v...), l.data)
}

func (l *boundLogger) Println(v ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogInfo, 2, fmt.Sprintln(v...), l.data)
}

func (l *boundLogger) GoFatal(v ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogFatal, 2, fmt.Sprint(v...), l.data)
	exitFunc(1)
}

func (l *boundLogger) GoFatalf(format string, v ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogFatal, 2, fmt.Sprintf(format, v...), l.data)
	exitFunc(1)
}

func (l *boundLogger) GoFatalln(v ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.l.log(time.Now(), LogFatal, 2, fmt.Sprintln(v...), l.data)
	exitFunc(1)
}

func (l *boundLogger) GoPanic(v ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	s := fmt.Sprint(v...)
	l.l.log(time.Now(), LogFatal, 2, s, l.data)
	panic(s)
}

func (l *boundLogger) GoPanicf(format string, v ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	s := fmt.Sprintf(format, v...)
	l.l.log(time.Now(), LogFatal, 2, s, l.data)
	panic(s)
}

func (l *boundLogger) GoPanicln(v ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	s := fmt.Sprintln(v...)
	l.l.log(time.Now(), LogFatal, 2, s, l.data)
	panic(s)
}

//-- LogBuilder interface -----------------

func (l *boundLogger) WithTimeFmt(format string) Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) ToStdout() Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) ToStderr() Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) ToWriter(out io.Writer) Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) WithPrefix(prefix string) Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) WithLine() Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) WithLn()  Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) ToFile(directory string, filename string) Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) ToAppender(appender Appender) Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) WithRotation(frequency rollFrequency, keepNLogs int) Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) WithStderr() Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) WithFmt(format string) Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) Json() Log5Go {
	// NOOP
	return l
}

func (l *boundLogger) Register(key string) (_ Log5Go, _ error) {
	// NOOP
	return l, fmt.Errorf("can't call register after calling WithData()")
}
