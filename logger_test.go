package log5go

import (
	"bytes"
	"fmt"
	"testing"
	"time"
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
