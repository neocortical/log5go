package log5go

import (
	"testing"
	"time"
)

func TestArchiveFilenamesDifferAcrossDST(t *testing.T) {
	// hourly across DST start
	t1, _ := time.Parse(time.RFC822, "09 Mar 14 01:00 PST")
	t2 := calculateNextRollTime(t1, RollHourly)

	f1 := generateArchiveFilename("foo.log", t1, RollHourly)
	f2 := generateArchiveFilename("foo.log", t2, RollHourly)

	if f1 == f2 {
		t.Errorf("filenames should differ across DST start")
	}

	// hourly across DST end
	t1, _ = time.Parse(time.RFC822, "02 Nov 14 01:00 PDT")
	t2 = calculateNextRollTime(t1, RollHourly)

	f1 = generateArchiveFilename("foo.log", t1, RollHourly)
	f2 = generateArchiveFilename("foo.log", t2, RollHourly)

	if f1 == f2 {
		t.Errorf("filenames should differ across DST end")
	}
}
