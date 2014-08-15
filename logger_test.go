package log5go

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
	"time"
)

const (
	Rxdate         					= `[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]`
	Rxtime         					= `[0-9][0-9]:[0-9][0-9]:[0-9][0-9]`
	Rxlevel									= `(TRACE|DEBUG|INFO|WARN|ERROR|FATAL|CUSTOM)`
	Rxprefix								= `prefix`
	Rxcaller								= `[a-zA-Z0-9_\-\/\.]+\.go`
	Rxline									= `[0-9]+`
	RxdefaultFmt						= Rxdate + " " + Rxtime + " " + Rxlevel + " : " + Rxmessage
	RxdefaultPrefixFmt 			= Rxdate + " " + Rxtime + " " + Rxlevel + " " + Rxprefix + ": " + Rxmessage
	RxdefaultLinesFmt 			= Rxdate + " " + Rxtime + " " + Rxlevel + ` \(` + Rxcaller + ":" + Rxline + `\): ` + Rxmessage
	RxdefaultPrefixLinesFmt = Rxdate + " " + Rxtime + " " + Rxlevel + " " + Rxprefix + ` \(` + Rxcaller + ":" + Rxline + `\): ` + Rxmessage
	Rxmessage								= `hello, world`
	Rxcustomformat 					= Rxmessage + ", " + Rxprefix + ", " + Rxlevel + "!!!"
)

func TestOutputOfMultipleLines(t *testing.T) {
	year := time.Now().Year()
	l := getTestLogger(LogAll)
	l.Trace("foo: %d", 1)
	l.Debug("bar: %d", 2)
	l.Info("baz: %d", 3)
	l.Warn("qux: %d", 4)
	l.Error("quux: %d", 5)
	l.Fatal("corge: %d", 6)

	a, _ := l.appender.(*bufferAppender)
	expected := fmt.Sprintf("%d TRACE : foo: 1\n%d DEBUG : bar: 2\n%d INFO : baz: 3\n%d WARN : qux: 4\n%d ERROR : quux: 5\n%d FATAL : corge: 6\n", year, year, year, year, year, year)
	if a.buf.String() != expected {
		t.Errorf("unexpected log output. expected \n%s\n ...but got \n%s", expected, a.buf.String())
	}
}

func TestDefaultFormats(t *testing.T) {
	var buf bytes.Buffer
	log := Logger(LogAll).ToWriter(&buf)
	runTest(log, &buf, RxdefaultFmt, t)

	log = Logger(LogAll).ToWriter(&buf).WithLn()
	runTest(log, &buf, RxdefaultLinesFmt, t)

	log = Logger(LogAll).ToWriter(&buf).WithPrefix("prefix")
	runTest(log, &buf, RxdefaultPrefixFmt, t)

	log = Logger(LogAll).ToWriter(&buf).WithPrefix("prefix").WithLine()
	runTest(log, &buf, RxdefaultPrefixLinesFmt, t)
}


func TestCustomFormat(t *testing.T) {
	var buf bytes.Buffer
	log := Logger(LogAll).ToWriter(&buf).WithPrefix("prefix").WithFmt("%m, %p, %l!!!")
	runTest(log, &buf, Rxcustomformat, t)
}

func runTest(log Log5Go, buf *bytes.Buffer, fmt string, t *testing.T) {
	buf.Reset()
	log.Info("hello, world")
	fmt = `^` + fmt + "\n$"
	matched, err := regexp.MatchString(fmt, buf.String())
	if err != nil || !matched {
		t.Errorf("expected \n%s but got \n%s", fmt, buf.String())
	}
}

func getTestLogger(level LogLevel) *logger {
	return &logger{
		level: level,
		appender: &bufferAppender{bytes.Buffer{}},
		timeFormat: "2006", // simple
	}
}

type bufferAppender struct {
	buf bytes.Buffer
}

func (a *bufferAppender) Append(msg string, level LogLevel, tstamp time.Time) {
	a.buf.Write([]byte(msg))
}
