package log5go

import (
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
