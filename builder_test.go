package log5go

import (
	"reflect"
	"testing"
)

func TestLog(t *testing.T) {
	lb := Log(LogAll)
	b, ok := lb.(*logBuilder)
	if !ok {
		t.Errorf("type of returned LogBuilder unexpected: %v", reflect.TypeOf(b))
	}

	if b.level != LogAll {
		t.Errorf("expected log level LogAll but was %d", b.level)
	}
	if b.appender != nil {
		t.Error("appender expected to be nil initially, but wasn't")
	}
	if b.timeFormat != TF_GoStd {
		t.Errorf("expected TF_GoStd default time format but was %d", b.timeFormat)
	}
	if b.errs.hasErrors() {
		t.Errorf("New logBuilder initialized with errors: %v", b.errs)
	}

	lb = Log(LogWarn)
	b, _ = lb.(*logBuilder)
	if b.level != LogWarn {
		t.Errorf("expected log level LogWarn but was %d", b.level)
	}
}

func TestWithTimeFmt(t *testing.T) {
	lb := Log(LogAll).WithTimeFmt(TF_NCSA)
	b, _ := lb.(*logBuilder)

	if b.timeFormat != TF_NCSA {
		t.Errorf("expected TF_GoStd after setting but was %d", b.timeFormat)
	}

	// custom time format, setting twice OK
	lb = Log(LogAll).WithTimeFmt(TF_NCSA).WithTimeFmt("2006,01,02")
	b, _ = lb.(*logBuilder)
	if b.timeFormat != "2006,01,02" {
		t.Errorf("expected '2006,01,02' after setting but was %d", b.timeFormat)
	}
}

func TestToStdout(t *testing.T) {
	lb := Log(LogAll).ToStdout()
	b, _ := lb.(*logBuilder)

	a, ok := b.appender.(*consoleAppender)
	if !ok {
		t.Errorf("Expected console appender but got %v", reflect.TypeOf(b.appender))
	}
	if a.stderrAware {
		t.Error("consoleAppender stderr expected false but was true")
	}

	// setting twice should result in error
	lb = Log(LogAll).ToStdout().ToStdout()
	b, _ = lb.(*logBuilder)
	if !b.errs.hasErrors() {
		t.Error("expected errors after setting appender twice but got none")
	}
}

func TestWithStdErr(t *testing.T) {
	lb := Log(LogAll).ToStdout().WithStderr()
	b, _ := lb.(*logBuilder)

	a, ok := b.appender.(*consoleAppender)
	if !ok {
		t.Errorf("Expected console appender but got %v", reflect.TypeOf(b.appender))
	}
	if !a.stderrAware {
		t.Error("consoleAppender stderr expected true but was false")
	}
}

// TODO: rest of builder methods
