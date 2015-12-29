package log5go

import (
	"os"
	"testing"
	"time"
)

func TestArchiveFilenamesDifferAcrossDST(t *testing.T) {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	// hourly across DST start
	t1, _ := time.ParseInLocation(time.RFC822, "09 Mar 14 01:00 PST", loc)
	t2 := calculateNextRollTime(t1, RollHourly)

	f1 := generateArchiveFilename("foo.log", t1, RollHourly)
	f2 := generateArchiveFilename("foo.log", t2, RollHourly)

	if f1 == f2 {
		t.Errorf("filenames should differ across DST start")
	}
	if f2 != "foo.log.2014-03-09-03-00-PDT" {
		t.Errorf("expected filename 'foo.log.2014-03-09-03-00-PDT' but got %s", f2)
	}

	// hourly across DST end
	t1, _ = time.ParseInLocation(time.RFC822, "02 Nov 14 01:00 PDT", loc)
	t2 = calculateNextRollTime(t1, RollHourly)

	f1 = generateArchiveFilename("foo.log", t1, RollHourly)
	f2 = generateArchiveFilename("foo.log", t2, RollHourly)

	if f1 == f2 {
		t.Errorf("filenames should differ across DST end")
	}
	if f2 != "foo.log.2014-11-02-01-00-PST" {
		t.Errorf("expected filename 'foo.log.2014-03-09-03-00-PDT' but got %s", f2)
	}
}

func TestWatchFiles(t *testing.T) {
	fname := "/tmp/reopentest.log"
	logfile, err := os.Create(fname)
	if err != nil {
		t.Errorf("error opening test file: %v", err)
	}
	a := &fileAppender{
		f:             logfile,
		fname:         fname,
		lastOpenTime:  time.Now(),
		nextRollTime:  time.Now(),
		rollFrequency: RollNone,
		keepNLogs:     SaveAllLogs,
	}
	fileAppenderMap[fname] = a

	if !a.validFile() {
		t.Error("expected validFile() to return true but got false")
	}

	// remove the file out from under the appender
	err = os.Remove(fname)
	if err != nil {
		t.Errorf("error deleting test file: %v", err)
	}

	if a.validFile() {
		t.Error("expected validFile() to return false but got true")
	}
}

func TestAutoReopen(t *testing.T) {
	fname := "/tmp/reopentest.log"
	logfile, err := os.Create(fname)
	if err != nil {
		t.Errorf("error opening test file: %v", err)
	}
	a := &fileAppender{
		f:             logfile,
		fname:         fname,
		lastOpenTime:  time.Now(),
		nextRollTime:  time.Now(),
		rollFrequency: RollNone,
		keepNLogs:     SaveAllLogs,
	}
	fileAppenderMap[fname] = a

	// remove the file out from under the appender
	err = os.Remove(fname)
	if err != nil {
		t.Errorf("error deleting test file: %v", err)
	}

	// detect failure
	watchFiles(time.Now())

	msg := []byte("hello")

	// should error: reopen refractory period not reached
	err = a.Append(&msg, LogInfo, time.Now())
	if err == nil {
		t.Error("should have errored writing to deleted file")
	}

	a.lastOpenTime = time.Now().Add(-time.Hour)

	watchFiles(time.Now())

	// should not error: file should be automatically reopened
	err = a.Append(&msg, LogInfo, time.Now())
	if err != nil {
		t.Errorf("failed to reopen test file: %v", err)
	}
}
