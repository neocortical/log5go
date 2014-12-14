package log5go

import (
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	l := Logger(LogWarn)
	ll, ok := l.(*logger)
	if !ok {
		t.Errorf("type of returned Log5Go unexpected: %v", reflect.TypeOf(l))
	}

	if ll.level != LogWarn {
		t.Errorf("expected log level LogAll but was %d", ll.level)
	}
	if ll.appender == nil {
		t.Error("appender should not be nil")
	}
	if ll.timeFormat != TF_GoStd {
		t.Errorf("expected TF_GoStd default time format but was %d", ll.timeFormat)
	}

	a, ok := ll.appender.(*writerAppender)
	if !ok {
		t.Errorf("Expected console appender but got %v", reflect.TypeOf(ll.appender))
	}
	if a.dest != os.Stderr {
		t.Errorf("writerAppender to stderr expected but was %v", a.dest)
	}
}

func TestWithTimeFmt(t *testing.T) {
	l := Logger(LogAll).WithTimeFmt(TF_NCSA)
	ll, _ := l.(*logger)

	if ll.timeFormat != TF_NCSA {
		t.Errorf("expected TF_GoStd after setting but was %d", ll.timeFormat)
	}

	// custom time format, setting twice OK
	l = Logger(LogAll).WithTimeFmt(TF_NCSA).WithTimeFmt("2006,01,02")
	ll, _ = l.(*logger)
	if ll.timeFormat != "2006,01,02" {
		t.Errorf("expected '2006,01,02' after setting but was %d", ll.timeFormat)
	}
}

func TestToStdout(t *testing.T) {
	l := Logger(LogAll).ToStdout()
	ll, _ := l.(*logger)

	a, ok := ll.appender.(*writerAppender)
	if !ok {
		t.Errorf("Expected console appender but got %v", reflect.TypeOf(ll.appender))
	}
	if a.dest != os.Stdout {
		t.Errorf("writerAppender to stdout expected but was %v", a.dest)
	}
	if a.errDest != nil {
		t.Errorf("writerAppender stderr split expected nil but was %v", a.errDest)
	}
}

func TestToStderr(t *testing.T) {
	l := Logger(LogAll).ToStderr()
	ll, _ := l.(*logger)

	a, ok := ll.appender.(*writerAppender)
	if !ok {
		t.Errorf("Expected writerAppender but got %v", reflect.TypeOf(ll.appender))
	}
	if a.dest != os.Stderr {
		t.Errorf("writerAppender to stderr expected but was %v", a.dest)
	}
	if a.errDest != nil {
		t.Errorf("writerAppender stderr split expected nil but was %v", a.errDest)
	}
}

func TestToWriter(t *testing.T) {
	out := new(bytes.Buffer)
	l := Logger(LogAll).ToWriter(out)
	ll, _ := l.(*logger)

	a, ok := ll.appender.(*writerAppender)
	if !ok {
		t.Errorf("Expected console appender but got %v", reflect.TypeOf(ll.appender))
	}
	if a.dest != out {
		t.Errorf("writerAppender to writer expected but was %v", a.dest)
	}
	if a.errDest != nil {
		t.Errorf("writerAppender stderr split expected nil but was %v", a.errDest)
	}
}

func TestWithStdErr(t *testing.T) {
	l := Logger(LogAll).ToStdout().WithStderr()
	ll, _ := l.(*logger)

	a, _ := ll.appender.(*writerAppender)
	if a.dest != os.Stdout {
		t.Errorf("writerAppender stdout expected but was %v", a.dest)
	}
	if a.errDest != os.Stderr {
		t.Errorf("writerAppender stderr expected but was %v", a.errDest)
	}
}

func TestToFile(t *testing.T) {
	l := Logger(LogAll).ToFile("/tmp", "foo.log")
	ll, _ := l.(*logger)

	a, ok := ll.appender.(*fileAppender)
	if !ok {
		t.Errorf("Expected fileAppender but got %v", reflect.TypeOf(ll.appender))
	}
	if a.f == nil {
		t.Error("fileAppender has nil file but should be initialized")
	}

	// initial rotation values
	if !time.Now().After(a.nextRollTime) {
		t.Errorf("expected nextRollTime to be Now() but was in the future")
	}
	if a.rollFrequency != RollNone {
		t.Errorf("expected rollFrequency init'd to none but was %d", a.rollFrequency)
	}
	if a.keepNLogs != SaveAllLogs {
		t.Errorf("expected keepNLogs to be SaveAllLogs but was %d", a.keepNLogs)
	}
}

func TestToFileMultipleAppenders(t *testing.T) {
	l := Logger(LogAll).ToStdout().ToFile("/tmp", "bar.log")
	ll, _ := l.(*logger)
	_, ok := ll.appender.(*fileAppender)
	if !ok {
		t.Error("expected fileAppender but was %v", reflect.TypeOf(ll.appender))
	}

	l = Logger(LogAll).ToFile("/tmp", "foo.log").ToStdout()
	ll, _ = l.(*logger)
	_, ok = ll.appender.(*writerAppender)
	if !ok {
		t.Error("expected writerAppender but was %v", reflect.TypeOf(ll.appender))
	}
}

func TestWithRotation(t *testing.T) {
	l := Logger(LogAll).ToFile("/tmp", "foo.log").WithRotation(RollDaily, 3)
	ll, _ := l.(*logger)
	a, _ := ll.appender.(*fileAppender)

	nextRoll := calculateNextRollTime(time.Now(), RollDaily)
	if !nextRoll.Equal(a.nextRollTime) {
		t.Errorf("expected nextRollTime to be %v but was %v", nextRoll, a.nextRollTime)
	}
	if a.rollFrequency != RollDaily {
		t.Errorf("expected rollFrequency %d but was %d", RollDaily, a.rollFrequency)
	}
	if a.keepNLogs != 3 {
		t.Errorf("expected keepNLogs to be 3 but was %d", a.keepNLogs)
	}

	// ToFile() must be called first
	l = Logger(LogAll).WithRotation(RollDaily, 3)
	ll, _ = l.(*logger)
	_, ok := ll.appender.(*writerAppender)
	if !ok {
		t.Error("expected writerAppender but was %v", reflect.TypeOf(ll.appender))
	}

	// can't apply to console logger
	l = Logger(LogAll).ToStdout().WithRotation(RollDaily, 3)
	ll, _ = l.(*logger)
	_, ok = ll.appender.(*writerAppender)
	if !ok {
		t.Error("expected writerAppender but was %v", reflect.TypeOf(ll.appender))
	}
}

func TestSameFileResultsInSameAppender(t *testing.T) {
	l1 := Logger(LogAll).ToFile("/tmp", "foo.log")
	ll1, _ := l1.(*logger)
	l2 := Logger(LogAll).ToFile("/tmp", "foo.log")
	ll2, _ := l2.(*logger)
	l3 := Logger(LogAll).ToFile("/tmp", "bar.log")
	ll3, _ := l3.(*logger)

	if ll1.appender == nil || ll2.appender == nil || ll3.appender == nil {
		t.Error("sanity check failed. expected all appenders to be init'd")
	}

	if ll1.appender == ll3.appender {
		t.Error("ll1 and ll3 appenders should not be the same object")
	}
	if ll1.appender != ll2.appender {
		t.Error("ll1 and ll2 appenders should be the same object")
	}
}

func TestRegister(t *testing.T) {
	log1 := Logger(LogInfo).ToStdout()
	log2 := Logger(LogInfo).ToStdout()

	if log1 == nil || log2 == nil {
		t.Error("sanity check failed, both logs should be init'd")
	}
	if log1 == log2 {
		t.Error("log1 and log2 should not be equal")
	}

	log3 := Logger(LogInfo).ToStdout().Register("foobar")
	log4, err := GetLog("foobar")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if log3 != log4 {
		t.Error("log3 and log4 should be equal but aren't")
	}
}

type nilAppender struct{}

func (a *nilAppender) Append(msg *[]byte, level LogLevel, tstamp time.Time) error {
	// NOOP
	return nil
}

func TestToAppender(t *testing.T) {
	appender := &nilAppender{}
	l := Logger(LogAll).ToAppender(appender)
	ll, _ := l.(*logger)

	if ll.appender == nil {
		t.Error("expected appender set after ToAppender() but is nil")
	}
	a, _ := ll.appender.(*nilAppender)
	if a != appender {
		t.Error("expected logger's appender to be equal to ToAppender() argument")
	}
}

func TestWithPrefix(t *testing.T) {
	l := Logger(LogAll).WithPrefix("foo")
	ll, _ := l.(*logger)

	if ll.prefix != "foo" {
		t.Errorf("expected prefix 'foo' but got %s", ll.prefix)
	}
}

func TestWithLines(t *testing.T) {
	l := Logger(LogAll)
	ll, _ := l.(*logger)
	if ll.lines != 0 {
		t.Errorf("expected default of no lines but got %d", ll.lines)
	}

	l = Logger(LogAll).WithShortLines()
	ll, _ = l.(*logger)
	if ll.lines != LogLinesShort {
		t.Errorf("expected short lines but got %d", ll.lines)
	}

	l = Logger(LogAll).WithLongLines()
	ll, _ = l.(*logger)
	if ll.lines != LogLinesLong {
		t.Errorf("expected long lines but got %d", ll.lines)
	}
}
