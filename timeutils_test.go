package log5go

import (
	"testing"
	"time"
)

const dateFmt = "2006-01-02T15:04:05"

func TestMinutelyOperations(t *testing.T) {
	traw, _ := time.Parse(dateFmt, "2014-08-11T12:15:22")
	t0, _ := time.Parse(dateFmt, "2014-08-11T12:15:00")
	t1, _ := time.Parse(dateFmt, "2014-08-11T12:16:00")
	t2, _ := time.Parse(dateFmt, "2014-08-11T12:17:00")

	tcalc := calculateNextRollTime(traw, RollMinutely)
	if !tcalc.Equal(t1) {
		t.Errorf("Error getting t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculateNextRollTime(tcalc, RollMinutely)
	if !tcalc.Equal(t2) {
		t.Errorf("Error getting t2. Expected %v but got %v", t2, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollMinutely)
	if !tcalc.Equal(t1) {
		t.Errorf("Error rolling back to t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollMinutely)
	if !tcalc.Equal(t0) {
		t.Errorf("Error getting t0. Expected %v but got %v", t0, tcalc)
	}
}

func TestHourlyOperations(t *testing.T) {
	traw, _ := time.Parse(dateFmt, "2014-08-11T23:15:22")
	t0, _ := time.Parse(dateFmt, "2014-08-11T23:00:00")
	t1, _ := time.Parse(dateFmt, "2014-08-12T00:00:00")
	t2, _ := time.Parse(dateFmt, "2014-08-12T01:00:00")

	tcalc := calculateNextRollTime(traw, RollHourly)
	if !tcalc.Equal(t1) {
		t.Errorf("Error getting t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculateNextRollTime(tcalc, RollHourly)
	if !tcalc.Equal(t2) {
		t.Errorf("Error getting t2. Expected %v but got %v", t2, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollHourly)
	if !tcalc.Equal(t1) {
		t.Errorf("Error rolling back to t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollHourly)
	if !tcalc.Equal(t0) {
		t.Errorf("Error getting t0. Expected %v but got %v", t0, tcalc)
	}
}

func TestDailyOperations(t *testing.T) {
	traw, _ := time.Parse(dateFmt, "2014-08-11T23:15:22")
	t0, _ := time.Parse(dateFmt, "2014-08-11T00:00:00")
	t1, _ := time.Parse(dateFmt, "2014-08-12T00:00:00")
	t2, _ := time.Parse(dateFmt, "2014-08-13T00:00:00")

	tcalc := calculateNextRollTime(traw, RollDaily)
	if !tcalc.Equal(t1) {
		t.Errorf("Error getting t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculateNextRollTime(tcalc, RollDaily)
	if !tcalc.Equal(t2) {
		t.Errorf("Error getting t2. Expected %v but got %v", t2, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollDaily)
	if !tcalc.Equal(t1) {
		t.Errorf("Error rolling back to t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollDaily)
	if !tcalc.Equal(t0) {
		t.Errorf("Error getting t0. Expected %v but got %v", t0, tcalc)
	}
}

func TestWeeklyOperations(t *testing.T) {
	traw, _ := time.Parse(dateFmt, "2014-08-11T23:15:22") // a Monday
	t0, _ := time.Parse(dateFmt, "2014-08-10T00:00:00")
	t1, _ := time.Parse(dateFmt, "2014-08-17T00:00:00")
	t2, _ := time.Parse(dateFmt, "2014-08-24T00:00:00")

	tcalc := calculateNextRollTime(traw, RollWeekly)
	if !tcalc.Equal(t1) {
		t.Errorf("Error getting t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculateNextRollTime(tcalc, RollWeekly)
	if !tcalc.Equal(t2) {
		t.Errorf("Error getting t2. Expected %v but got %v", t2, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollWeekly)
	if !tcalc.Equal(t1) {
		t.Errorf("Error rolling back to t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollWeekly)
	if !tcalc.Equal(t0) {
		t.Errorf("Error getting t0. Expected %v but got %v", t0, tcalc)
	}
}

func TestHourlyHandlesDSTStart(t *testing.T) {

	// err is ok here. times will be UTC and code/test will behave correctly
	loc, _ := time.LoadLocation("America/Los_Angeles")

	t0, _ := time.ParseInLocation(time.RFC822, "09 Mar 14 00:15 PST", loc)
	t1, _ := time.ParseInLocation(time.RFC822, "09 Mar 14 01:00 PST", loc)
	t3, _ := time.ParseInLocation(time.RFC822, "09 Mar 14 03:00 PDT", loc)

	tnext := calculateNextRollTime(t0, RollHourly)
	if !tnext.Equal(t1) {
		t.Errorf("Error getting t1. Expected %v but got %v", t1, tnext)
	}

	tnext = calculateNextRollTime(tnext, RollHourly)
	if !tnext.Equal(t3) {
		t.Errorf("Error getting t3. Expected %v but got %v", t3, tnext)
	}

	tprev := calculatePreviousRollTime(tnext, RollHourly)
	if !tprev.Equal(t1) {
		t.Errorf("Error rolling back to t1. Expected %v but got %v", t1, tprev)
	}
}

func TestHourlyHandlesDSTEnd(t *testing.T) {

	// err is ok here. times will be UTC and code/test will behave correctly
	loc, _ := time.LoadLocation("America/Los_Angeles")

	t0, _ := time.ParseInLocation(time.RFC822, "02 Nov 14 00:15 PDT", loc)
	t1, _ := time.ParseInLocation(time.RFC822, "02 Nov 14 01:00 PDT", loc)
	t2, _ := time.ParseInLocation(time.RFC822, "02 Nov 14 01:00 PST", loc)

	tnext := calculateNextRollTime(t0, RollHourly)
	if !tnext.Equal(t1) {
		t.Errorf("Error getting t1. Expected %v but got %v", t1, tnext)
	}

	tnext = calculateNextRollTime(tnext, RollHourly)
	if !tnext.Equal(t2) {
		t.Errorf("Error getting t2. Expected %v but got %v", t2, tnext)
	}

	tprev := calculatePreviousRollTime(tnext, RollHourly)
	if !tprev.Equal(t1) {
		t.Errorf("Error rolling back to t1. Expected %v but got %v", t1, tprev)
	}
}

func TestDailyHandlesDST(t *testing.T) {
	// DST start
	traw, _ := time.Parse(time.RFC822, "09 Mar 14 12:15 PDT")
	t0, _ := time.Parse(time.RFC822, "09 Mar 14 00:00 PST")
	t1, _ := time.Parse(time.RFC822, "10 Mar 14 00:00 PDT")

	tcalc := calculateNextRollTime(traw, RollDaily)
	if !tcalc.Equal(t1) {
		t.Errorf("Error getting t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollDaily)
	if !tcalc.Equal(t0) {
		t.Errorf("Error rolling back to t0. Expected %v but got %v", t0, tcalc)
	}

	// DST end
	traw, _ = time.Parse(time.RFC822, "02 Nov 14 12:15 PST")
	t0, _ = time.Parse(time.RFC822, "02 Nov 14 00:00 PDT")
	t1, _ = time.Parse(time.RFC822, "03 Nov 14 00:00 PST")

	tcalc = calculateNextRollTime(traw, RollDaily)
	if !tcalc.Equal(t1) {
		t.Errorf("Error getting t1. Expected %v but got %v", t1, tcalc)
	}

	tcalc = calculatePreviousRollTime(tcalc, RollDaily)
	if !tcalc.Equal(t0) {
		t.Errorf("Error rolling back to t0. Expected %v but got %v", t0, tcalc)
	}
}
