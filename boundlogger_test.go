package log5go

import (
	"bytes"
	"os"
	"testing"
)

var boundLoggerTests = []loggerTest{
	{"withdata", "hello, %s", []interface{}{"world"}, Data{"foo": "bar"}, TF_GoStd, "prefix", nil, "^" + Rxdate + " " + Rxtime + " {{level}} prefix: hello, world foo=\"bar\"\n$"},
}

func TestWithData(t *testing.T) {
	inner, _ := Logger(LogAll).(*logger)
	bl := &boundLogger{l: inner, data: make(Data)}

	if len(bl.data) != 0 {
		t.Errorf("expected empty data but got %v", bl.data)
	}

	bl, ok := bl.WithData(Data{"foo": "bar", "baz": 1}).(*boundLogger)
	if !ok {
		t.Error("expected *boundLogger back from WithData")
	}
	if len(bl.data) != 2 || bl.data["foo"] != "bar" || bl.data["baz"] != 1 {
		t.Errorf("expected two-element data but got %v", bl.data)
	}

	// piling on
	bl, ok = bl.WithData(Data{"baz": 2, "qux": 3.14}).(*boundLogger)
	if !ok {
		t.Error("expected *boundLogger back from WithData")
	}
	if len(bl.data) != 3 || bl.data["foo"] != "bar" || bl.data["baz"] != 2 || bl.data["qux"] != 3.14 {
		t.Errorf("expected three-element data but got %v", bl.data)
	}
}

func TestAllBoundLoggerLevels(t *testing.T) {
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
		bl := &boundLogger{
			l:    l,
			data: test.data,
		}
		runLoggerLevelTest(bl, &buf, &test, t)
	}

	for _, test := range boundLoggerTests {
		l := &logger{
			level:      LogAll,
			formatter:  test.formatter,
			appender:   appender,
			timeFormat: test.timeFormat,
			prefix:     test.prefix,
		}
		bl := &boundLogger{
			l:    l,
			data: test.data,
		}
		runLoggerLevelTest(bl, &buf, &test, t)
	}
}

func TestBoundLoggerBuilderNoops(t *testing.T) {
	var buf bytes.Buffer
	appender := &writerAppender{dest: &buf}
	formatter := NewStringFormatter("%m")
	l := &logger{
		level:      LogAll,
		formatter:  formatter,
		appender:   appender,
		timeFormat: TF_GoStd,
		prefix:     "foo",
	}
	bl := &boundLogger{
		l:    l,
		data: Data{},
	}

	bl.SetOutput(os.Stdout)
	if l.appender != appender {
		t.Error("appender changed")
	}

	bl.SetFlags(1)
	if l.flag != 0 {
		t.Error("appender changed")
	}

	bl.SetPrefix("bar")
	if l.prefix != "foo" {
		t.Error("prefix changed")
	}

	bl2 := bl.WithTimeFmt("2006")
	if bl2 != bl || l.timeFormat != TF_GoStd {
		t.Error("time format changed")
	}

	bl2 = bl.ToStdout()
	if bl2 != bl || l.appender != appender {
		t.Error("appender changed")
	}

	bl2 = bl.ToStderr()
	if bl2 != bl || l.appender != appender {
		t.Error("appender changed")
	}

	bl2 = bl.ToWriter(os.Stdout)
	if bl2 != bl || l.appender != appender {
		t.Error("appender changed")
	}

	bl2 = bl.WithPrefix("bar")
	if bl2 != bl || l.prefix != "foo" {
		t.Error("prefix changed")
	}

	bl2 = bl.WithLine()
	if bl2 != bl || l.lines != 0 || l.flag != 0 {
		t.Error("lines changed")
	}

	bl2 = bl.WithLn()
	if bl2 != bl || l.lines != 0 || l.flag != 0 {
		t.Error("lines changed")
	}

	bl2 = bl.ToFile("/tmp", "foo.log")
	if bl2 != bl || l.appender != appender {
		t.Error("appender changed")
	}

	bl2 = bl.ToAppender(&writerAppender{dest: os.Stdout})
	if bl2 != bl || l.appender != appender {
		t.Error("appender changed")
	}

	bl2 = bl.ToFile("/tmp", "foo.log").WithRotation(RollDaily, 7)
	if bl2 != bl || l.appender != appender {
		t.Error("appender changed")
	}

	bl2 = bl.ToStdout().WithStderr()
	if bl2 != bl || l.appender != appender {
		t.Error("appender changed")
	}

	bl2 = bl.WithFmt("%t")
	if bl2 != bl || l.formatter != formatter {
		t.Error("formatter changed")
	}
	if len(formatter.parts) != 1 || formatter.parts[0] != "%m" {
		t.Error("formatter changed")
	}

	bl2 = bl.Json()
	if bl2 != bl || l.appender != appender {
		t.Error("appender changed")
	}
}

func TestBoundLoggerFatals(t *testing.T) {
	exitCalled := 0
	exitFunc = func(i int) {
		exitCalled = i
	}

	inner := &logger{
		level:      LogAll,
		formatter:  nil,
		appender:   &writerAppender{dest: os.Stdout},
		timeFormat: TF_GoStd,
		prefix:     "",
	}
	l := &boundLogger{
		l:    inner,
		data: Data{},
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

func TestBoundLoggerPanics(t *testing.T) {
	inner := &logger{
		level:      LogAll,
		formatter:  nil,
		appender:   &writerAppender{dest: os.Stdout},
		timeFormat: TF_GoStd,
		prefix:     "",
	}
	l := &boundLogger{
		l:    inner,
		data: Data{},
	}

	f := func(l *boundLogger, t *testing.T) {
		defer assertRecoverPanic(t)
		l.GoPanic("aiyeee!")
	}
	f(l, t)

	f = func(l *boundLogger, t *testing.T) {
		defer assertRecoverPanic(t)
		l.GoPanicf("[wilhelm scream]!")
	}
	f(l, t)

	f = func(l *boundLogger, t *testing.T) {
		defer assertRecoverPanic(t)
		l.GoPanicln("where am i?!")
	}
	f(l, t)
}
