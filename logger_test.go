package log5go

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type loggerTest struct {
	msg      string
	args     []interface{}
	data     Data
	expected string
	create   loggerFunc
}

type loggerFunc func() Log5Go

const (
	Rxdate                  = `[0-9]{4}/[0-9]{2}/[0-9]{2}`
	Rxtime                  = `[0-9]{2}:[0-9]{2}:[0-9]{2}`
	Rxlevel                 = `(TRACE|DEBUG|INFO|WARN|ERROR|FATAL|CUSTOM)`
	Rxprefix                = `prefix`
	Rxcaller                = `[a-zA-Z0-9_\-\/\.]+\.go`
	Rxline                  = `[0-9]+`
	RxdefaultFmt            = Rxdate + " " + Rxtime + " " + Rxlevel + " : " + Rxmessage
	RxdefaultPrefixFmt      = Rxdate + " " + Rxtime + " " + Rxlevel + " " + Rxprefix + ": " + Rxmessage
	RxdefaultLinesFmt       = Rxdate + " " + Rxtime + " " + Rxlevel + ` \(` + Rxcaller + ":" + Rxline + `\): ` + Rxmessage
	RxdefaultPrefixLinesFmt = Rxdate + " " + Rxtime + " " + Rxlevel + " " + Rxprefix + ` \(` + Rxcaller + ":" + Rxline + `\): ` + Rxmessage
	Rxmessage               = `hello, world`
	Rxcustomformat          = Rxmessage + ", " + Rxprefix + ", " + Rxlevel + "!!!"
	Rxdata                  = `(pi=3\.14159265359 foo=\"bar\"|foo=\"bar\" pi=3\.14159265359)`
)

var loggerTests = []loggerTest{
	{
		msg:      "hello",
		expected: "^" + Rxdate + " " + Rxtime + " {{level}} : hello\n$",
		create:   func() Log5Go { return Logger(LogAll) },
	},
	{
		msg:      "侍 (%s)",
		args:     []interface{}{"samurai"},
		expected: "^" + Rxdate + " " + Rxtime + " {{level}} : 侍 \\(samurai\\)\n$",
		create:   func() Log5Go { return Logger(LogAll) },
	},
	{
		msg:      "foo",
		args:     nil,
		expected: "^" + Rxdate + " " + Rxtime + " {{level}} : foo (foo=\"bar\"|pi=3.14) (foo=\"bar\"|pi=3.14)\n$",
		data:     Data{"foo": "bar", "pi": 3.14},
		create:   func() Log5Go { return Logger(LogAll) },
	},
}

func Test_RunLoggerTests(t *testing.T) {
	var buf bytes.Buffer
	appender := &writerAppender{dest: &buf}

	for _, test := range loggerTests {
		l := test.create()
		l = l.ToAppender(appender)

		if test.data != nil {
			l = l.WithData(test.data)
		}

		runLevelTest(t, test, l, LogTrace, &buf)
		runLevelTest(t, test, l, LogDebug, &buf)
		runLevelTest(t, test, l, LogInfo, &buf)
		runLevelTest(t, test, l, LogNotice, &buf)
		runLevelTest(t, test, l, LogWarn, &buf)
		runLevelTest(t, test, l, LogError, &buf)
		runLevelTest(t, test, l, LogCritical, &buf)
		runLevelTest(t, test, l, LogAlert, &buf)
		runLevelTest(t, test, l, LogFatal, &buf)
	}
}

func runLevelTest(t *testing.T, test loggerTest, l Log5Go, level LogLevel, buf *bytes.Buffer) {
	var expected string

	if test.args != nil {
		l.Log(level, test.msg, test.args...)
	} else {
		l.Log(level, test.msg)
	}
	expected = subLevel(test.expected, GetLogLevelString(level))
	buf.Reset()

	switch level {
	case LogTrace:
		if test.args != nil {
			l.Trace(test.msg, test.args...)
		} else {
			l.Trace(test.msg)
		}
		expected = subLevel(test.expected, "TRACE")
	case LogDebug:
		if test.args != nil {
			l.Debug(test.msg, test.args...)
		} else {
			l.Debug(test.msg)
		}
		expected = subLevel(test.expected, "DEBUG")
	case LogInfo:
		if test.args != nil {
			l.Info(test.msg, test.args...)
		} else {
			l.Info(test.msg)
		}
		expected = subLevel(test.expected, "INFO")
	case LogNotice:
		if test.args != nil {
			l.Notice(test.msg, test.args...)
		} else {
			l.Notice(test.msg)
		}
		expected = subLevel(test.expected, "NOTICE")
	case LogWarn:
		if test.args != nil {
			l.Warn(test.msg, test.args...)
		} else {
			l.Warn(test.msg)
		}
		expected = subLevel(test.expected, "WARN")
	case LogError:
		if test.args != nil {
			l.Error(test.msg, test.args...)
		} else {
			l.Error(test.msg)
		}
		expected = subLevel(test.expected, "ERROR")
	case LogCritical:
		if test.args != nil {
			l.Critical(test.msg, test.args...)
		} else {
			l.Critical(test.msg)
		}
		expected = subLevel(test.expected, "CRIT")
	case LogAlert:
		if test.args != nil {
			l.Alert(test.msg, test.args...)
		} else {
			l.Alert(test.msg)
		}
		expected = subLevel(test.expected, "ALERT")
	case LogFatal:
		if test.args != nil {
			l.Fatal(test.msg, test.args...)
		} else {
			l.Fatal(test.msg)
		}
		expected = subLevel(test.expected, "FATAL")

	}
	assertMatch(t, expected, buf)

	buf.Reset()
}

func assertMatch(t *testing.T, expected string, buf *bytes.Buffer) {
	matched, err := regexp.MatchString(expected, buf.String())
	assert.Nil(t, err, "regexp error matching output: %v", err)
	assert.True(t, matched, "expected \n%s\nbut got\n%s\n", expected, buf.String())
}

func subLevel(expected, level string) string {
	return strings.Replace(expected, "{{level}}", level, -1)
}

func TestOutputOfMultipleLines(t *testing.T) {
	year := time.Now().Year()
	l := Logger(LogAll).WithTimeFmt("2006").ToAppender(&bufferAppender{bytes.Buffer{}}).(*logger)
	l.Trace("foo: %d", 1)
	l.Debug("bar: %d", 2)
	l.Info("baz: %d", 3)
	l.Warn("qux: %d", 4)
	l.Error("quux: %d", 5)
	l.Fatal("corge: %d", 6)

	a := l.appender.(*bufferAppender)
	expected := fmt.Sprintf("%d TRACE : foo: 1\n%d DEBUG : bar: 2\n%d INFO : baz: 3\n%d WARN : qux: 4\n%d ERROR : quux: 5\n%d FATAL : corge: 6\n", year, year, year, year, year, year)
	if a.buf.String() != expected {
		t.Errorf("unexpected log output. expected \n%s\n ...but got \n%s", expected, a.buf.String())
	}
}

func TestDefaultFormats(t *testing.T) {
	var buf bytes.Buffer
	log := Logger(LogAll).ToWriter(&buf)
	runTest(log, &buf, RxdefaultFmt, t)

	log = Logger(LogAll).ToWriter(&buf).WithShortLines()
	runTest(log, &buf, RxdefaultLinesFmt, t)

	log = Logger(LogAll).ToWriter(&buf).WithPrefix("prefix")
	runTest(log, &buf, RxdefaultPrefixFmt, t)

	log = Logger(LogAll).ToWriter(&buf).WithPrefix("prefix").WithLongLines()
	runTest(log, &buf, RxdefaultPrefixLinesFmt, t)
}

func TestCustomFormat(t *testing.T) {
	var buf bytes.Buffer
	log := Logger(LogAll).ToWriter(&buf).WithPrefix("prefix").WithFmt("%m, %p, %l!!!")
	runTest(log, &buf, Rxcustomformat, t)
}

func TestDataStringFormatter(t *testing.T) {
	var buf bytes.Buffer
	log := Logger(LogAll).ToWriter(&buf).WithFmt("%m")

	runTest(log.WithData(Data{"foo": "bar", "pi": 3.14159265359}), &buf, Rxmessage+" "+Rxdata, t)
}

func TestScrubData(t *testing.T) {
	x := 1
	var badbuf bytes.Buffer
	badMap := map[int]string{1: "hi"}
	var okiface interface{} = reflect.ValueOf(x).Interface()
	var badiface interface{} = reflect.ValueOf(badbuf).Interface()
	var slice []byte
	var strct struct{}
	d := Data{"badMap": badMap, "okiface": okiface, "badiface": badiface, "slice": slice, "strct": strct, "bar": "baz"}

	d = scrubData(d)

	if len(d) != 2 || d["bar"] != "baz" || d["okiface"] != 1 {
		t.Errorf("expected single valid element but got: %v", d)
	}
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

type bufferAppender struct {
	buf bytes.Buffer
}

func (a *bufferAppender) Append(msg *[]byte, level LogLevel, tstamp time.Time) (err error) {
	TerminateMessageWithNewline(msg)
	_, err = a.buf.Write(*msg)
	return err
}

func TestGetSetLevel(t *testing.T) {
	var buf bytes.Buffer
	appender := &writerAppender{dest: &buf}
	l := &logger{
		level:      LogAll,
		formatter:  nil,
		appender:   appender,
		timeFormat: TF_GoStd,
		prefix:     "foo",
	}

	l.SetLogLevel(LogWarn)
	if l.LogLevel() != LogWarn {
		t.Errorf("expected %d but got %d", LogWarn, l.LogLevel())
	}
	if l.level != LogWarn {
		t.Errorf("expected %d but got %d", LogWarn, l.level)
	}
}
