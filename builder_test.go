package log5go

import (
	"os"
	"reflect"
	"testing"
	"time"
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
	if a.errDest != nil {
		t.Errorf("consoleAppender stderr expected nil but was %v", a.errDest)
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
	if a.errDest != os.Stderr {
		t.Errorf("consoleAppender stderr expected but was %v", a.errDest)
	}
}

func TestToFile(t *testing.T) {
	lb := Log(LogAll).ToFile("/tmp", "foo.log")
	b, _ := lb.(*logBuilder)

	a, ok := b.appender.(*fileAppender)
	if !ok {
		t.Errorf("Expected fileAppender but got %v", reflect.TypeOf(b.appender))
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
	if a.keepNLogs != SaveAllOldLogs {
		t.Errorf("expected keepNLogs to be SaveAllOldLogs but was %d", a.keepNLogs)
	}
}

func TestToFileMultipleAppenders(t *testing.T) {
	lb := Log(LogAll).ToFile("/tmp", "foo.log").ToFile("/tmp", "bar.log")
	b, _ := lb.(*logBuilder)
	if !b.errs.hasErrors() {
		t.Error("expected error after calling ToFile() twice but got none")
	}

	lb = Log(LogAll).ToFile("/tmp", "foo.log").ToStdout()
	b, _ = lb.(*logBuilder)
	if !b.errs.hasErrors() {
		t.Error("expected error after calling ToFile() and then ToConsole() but got none")
	}
}

func TestWithRotation(t *testing.T) {
	lb := Log(LogAll).ToFile("/tmp", "foo.log").WithRotation(RollDaily, 3)
	b, _ := lb.(*logBuilder)
	a, _ := b.appender.(*fileAppender)

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
	lb = Log(LogAll).WithRotation(RollDaily, 3)
	b, _ = lb.(*logBuilder)
	if !b.errs.hasErrors() {
		t.Error("expected error after illegal WithRotation() call but got none")
	}

	// can't apply to console logger
	lb = Log(LogAll).ToStdout().WithRotation(RollDaily, 3)
	b, _ = lb.(*logBuilder)
	if !b.errs.hasErrors() {
		t.Error("expected error after illegal WithRotation() call but got none")
	}
}

func TestSameFileResultsInSameAppender(t *testing.T) {
	lb1 := Log(LogAll).ToFile("/tmp", "foo.log")
	b1, _ := lb1.(*logBuilder)
	lb2 := Log(LogAll).ToFile("/tmp", "foo.log")
	b2, _ := lb2.(*logBuilder)
	lb3 := Log(LogAll).ToFile("/tmp", "bar.log")
	b3, _ := lb3.(*logBuilder)

	if b1.appender == nil || b2.appender == nil || b3.appender == nil {
		t.Error("sanity check failed. expected all appenders to be init'd")
	}

	if b1.appender == b3.appender {
		t.Error("b1 and b3 appenders should not be the same object")
	}
	if b1.appender != b2.appender {
		t.Error("b1 and b2 appenders should be the same object")
	}
}

func TestBuild(t *testing.T) {
	log, err := Log(LogInfo).ToStdout().Build()
	if err != nil {
		t.Errorf("expected no errors but got %v", err)
	}
	l, _ := log.(*logger)
	if l.level != LogInfo {
		t.Errorf("expected log level %d but got %d", LogInfo, l.level)
	}
	a, ok := l.appender.(*consoleAppender)
	if !ok {
		t.Error("built logger's appender should be console but isn't")
	}
	if a.errDest != nil {
		t.Error("console logger should not split errors by default")
	}
	if l.timeFormat != TF_GoStd {
		t.Error("expected default TF %d but got %d", TF_GoStd, l.timeFormat)
	}

	log, err = Log(LogError).ToFile("/tmp", "foo.log").Build()
	if err != nil {
		t.Errorf("expected no errors but got %v", err)
	}
	l, _ = log.(*logger)
	if l.level != LogError {
		t.Errorf("expected log level %d but got %d", LogError, l.level)
	}
	a2, ok := l.appender.(*fileAppender)
	if !ok {
		t.Error("built logger's appender should be file but isn't")
	}
	if a2.f == nil {
		t.Error("appender's file reference in nil")
	}
	if a2.f.Name() != "/tmp/foo.log" {
		t.Errorf("expected filename '/tmp/foo.log' but got %s", a2.f.Name())
	}

	// some multi-method sanity checks...
	log, err = Log(LogInfo).ToStdout().WithTimeFmt(TF_NCSA).WithStderr().Build()
	if err != nil {
		t.Errorf("error building logger with console, time, stderr: %v", err)
	}

	log, err = Log(LogInfo).WithTimeFmt(TF_NCSA).ToFile("/tmp/", "foo.log").WithRotation(RollMinutely, 10).Build()
	if err != nil {
		t.Errorf("error building logger with file, time, rotation: %v", err)
	}
}

func TestRegister(t *testing.T) {
	log1, _ := Log(LogInfo).ToStdout().Build()
	log2, _ := Log(LogInfo).ToStdout().Build()

	if log1 == nil || log2 == nil {
		t.Error("sanity check failed, both logs should be init'd")
	}
	if log1 == log2 {
		t.Error("log1 and log2 should not be equal")
	}

	log3, _ := Log(LogInfo).ToStdout().Register("foobar")
	log4, err := GetLog("foobar")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if log3 != log4 {
		t.Error("log3 and log4 should be equal but aren't")
	}

	_, err = Log(LogInfo).ToStdout().Register("foobar")
	if err == nil {
		t.Error("expected error trying to register name twice but got none")
	}
}

type nilAppender struct{}

func (a *nilAppender) Append(msg string, level LogLevel, tstamp time.Time) {
	// NOOP
}

func TestToAppender(t *testing.T) {
	appender := &nilAppender{}
	lb := Log(LogAll).ToAppender(appender)
	b, _ := lb.(*logBuilder)

	if b.errs.hasErrors() {
		t.Errorf("expected no errors after ToAppender() but got %v", b.errs.Error())
	}
	if b.appender == nil {
		t.Error("expected appender set after ToAppender() but is nil")
	}
	a, _ := b.appender.(*nilAppender)
	if a != appender {
		t.Error("expected builder's appender to be equal to ToAppender() argument")
	}
}
