package log5go

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

type loggerTest struct {
	testname   string
	msg        string
	args       []interface{}
	data       Data
	timeFormat string
	prefix     string
	formatter  Formatter
	expected   string
}

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
	{"basic", "hello", []interface{}{}, Data{}, "2006", "noshow", NewStringFormatter("%t %l %m"), "^[0-9]{4} {{level}} hello\n$"},
	{"stdwithutf8", "侍 (%s)", []interface{}{"samurai"}, Data{}, TF_GoStd, "prefix", nil, "^" + Rxdate + " " + Rxtime + " {{level}} prefix: 侍 \\(samurai\\)\n$"},
}

func TestAllLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	appender := &writerAppender{dest: &buf}

	for _, test := range loggerTests {
		l := &logger{
			level:      LogAll,
			formatter:  test.formatter,
			appender:   appender,
			timeFormat: test.timeFormat,
			prefix:     test.prefix,
		}
		runLoggerLevelTest(l, &buf, &test, t)
	}
}

func runLoggerLevelTest(l Log5Go, buf *bytes.Buffer, test *loggerTest, t *testing.T) {

	var LLCustom LogLevel = LogInfo + 1
	RegisterLogLevel(LLCustom, "CUSTOM")

	l.SetLogLevel(LogTrace)
	l.Trace(test.msg, test.args...)
	expected := subLevel(test.expected, "TRACE")
	assertMatch(test.testname, "trace", expected, buf, t)
	l.SetLogLevel(LogTrace + 1)
	l.Trace(test.msg, test.args...)
	assertMatch(test.testname, "tracethresh", "", buf, t)

	l.SetLogLevel(LogDebug)
	l.Debug(test.msg, test.args...)
	expected = subLevel(test.expected, "DEBUG")
	assertMatch(test.testname, "debug", expected, buf, t)
	l.SetLogLevel(LogDebug + 1)
	l.Debug(test.msg, test.args...)
	assertMatch(test.testname, "debugthresh", "", buf, t)

	l.SetLogLevel(LogAll)
	l.Info(test.msg, test.args...)
	expected = subLevel(test.expected, "INFO")
	assertMatch(test.testname, "info", expected, buf, t)
	l.Printf(test.msg, test.args...)
	assertMatch(test.testname, "printf", expected, buf, t)
	l.SetLogLevel(LogInfo + 1)
	l.Info(test.msg, test.args...)
	assertMatch(test.testname, "infothresh", "", buf, t)

	l.Warn(test.msg, test.args...)
	expected = subLevel(test.expected, "WARN")
	assertMatch(test.testname, "warn", expected, buf, t)
	l.SetLogLevel(LogWarn + 1)
	l.Warn(test.msg, test.args...)
	assertMatch(test.testname, "warnthresh", "", buf, t)

	l.Error(test.msg, test.args...)
	expected = subLevel(test.expected, "ERROR")
	assertMatch(test.testname, "error", expected, buf, t)
	l.SetLogLevel(LogError + 1)
	l.Error(test.msg, test.args...)
	assertMatch(test.testname, "errorthresh", "", buf, t)

	exitCalled := ""
	exitFunc = func(i int) {
		exitCalled = fmt.Sprintf("%d", i)
	}

	l.Fatal(test.msg, test.args...)
	if exitCalled != "" {
		t.Errorf("exit shouldn't have been called but was with %s", exitCalled)
	}
	expected = subLevel(test.expected, "FATAL")
	assertMatch(test.testname, "fatal", expected, buf, t)
	l.SetLogLevel(LogFatal + 1)
	l.Fatal(test.msg, test.args...)
	assertMatch(test.testname, "fatalthresh", "", buf, t)

	// test gofatalf
	l.SetLogLevel(LogAll)
	l.GoFatalf(test.msg, test.args...)
	if exitCalled != "1" {
		t.Errorf("exit not called or wrong output code given: %s", exitCalled)
	}
	assertMatch(test.testname, "gofatalf", "", buf, t)

	// cleanup
	exitFunc = os.Exit

	l.SetLogLevel(LLCustom)
	l.Log(LLCustom, test.msg, test.args...)
	expected = subLevel(test.expected, "CUSTOM")
	assertMatch(test.testname, "custom", expected, buf, t)
	l.SetLogLevel(LogWarn)
	l.Log(LLCustom, test.msg, test.args...)
	assertMatch(test.testname, "customthresh", "", buf, t)

	// cleanup
	DeregisterLogLevel(LLCustom)
}

func assertMatch(testname string, testpart, expected string, buf *bytes.Buffer, t *testing.T) {
	matched, err := regexp.MatchString(expected, buf.String())
	if err != nil || !matched {
		t.Errorf("test %s/%s expected \n%s but got \n%s", testname, testpart, expected, buf.String())
	}
	buf.Reset()
}

func subLevel(expected, level string) string {
	return strings.Replace(expected, "{{level}}", level, -1)
}

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

func getTestLogger(level LogLevel) *logger {
	return &logger{
		level:      level,
		appender:   &bufferAppender{bytes.Buffer{}},
		timeFormat: "2006", // simple
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

func TestLoggerFatals(t *testing.T) {
	exitCalled := 0
	exitFunc = func(i int) {
		exitCalled = i
	}

	l := &logger{
		level:      LogAll,
		formatter:  nil,
		appender:   &writerAppender{dest: os.Stdout},
		timeFormat: TF_GoStd,
		prefix:     "",
	}

	l.GoFatal("jeepers!")
	if exitCalled != 1 {
		t.Error("expected exit called, but wasn't")
	}

	exitCalled = 0
	l.GoFatalf("yoinks!")
	if exitCalled != 1 {
		t.Error("expected exit called, but wasn't")
	}

	exitCalled = 0
	l.GoFatalln("aye carumba!")
	if exitCalled != 1 {
		t.Error("expected exit called, but wasn't")
	}

	// cleanup
	exitFunc = os.Exit
}

func TestLoggerPanics(t *testing.T) {
	l := &logger{
		level:      LogAll,
		formatter:  nil,
		appender:   &writerAppender{dest: os.Stdout},
		timeFormat: TF_GoStd,
		prefix:     "",
	}

	f := func(l *logger, t *testing.T) {
		defer assertRecoverPanic(t)
		l.GoPanic("aiyeee!")
	}
	f(l, t)

	f = func(l *logger, t *testing.T) {
		defer assertRecoverPanic(t)
		l.GoPanicf("[wilhelm scream]!")
	}
	f(l, t)

	f = func(l *logger, t *testing.T) {
		defer assertRecoverPanic(t)
		l.GoPanicln("where am i?!")
	}
	f(l, t)
}

func assertRecoverPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Error("expected panic. got nothing.")
	}
}
