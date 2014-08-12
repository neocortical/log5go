package log5go

import (
	"time"
)

// given a time, calculate the instant that the log should next roll
func calculateNextRollTime(t time.Time, freq rollFrequency) time.Time {
	if freq == RollMinutely {
		t = t.Truncate(time.Minute)
		return t.Add(time.Minute)
	} else if freq == RollHourly {
		t = t.Truncate(time.Hour)
		t2 := t.Add(time.Hour)
		// daylight savings end test
		if t2.Hour() == t.Hour() {
			t2 = t2.Add(time.Hour)
		}
		return t2
	} else {
		t = t.Truncate(time.Hour)
		// easiest way to beat DST bugs is to just iterate
		for t.Hour() > 0 {
			t = t.Add(-time.Hour)
		}
		if freq == RollDaily {
			return t.AddDate(0, 0, 1)
		} else {
			if t.Weekday() == time.Sunday {
				return t.AddDate(0, 0, 7)
			}
			for t.Weekday() != time.Sunday {
				t = t.AddDate(0, 0, 1)
			}
			return t
		}
	}
}

// find the previous roll time for the given frequency.
// time t is assumed to be truncated to a valid roll time.
func calculatePreviousRollTime(t time.Time, freq rollFrequency) time.Time {
	switch freq {
	case RollMinutely:
		return t.Add(-time.Minute)
	case RollHourly:
		t := t.Add(-time.Hour)
		t2 := t.Add(-time.Hour)
		// daylight savings end test
		if t2.Hour() == t.Hour() {
			return t2
		} else {
			return t
		}
	case RollDaily:
		return t.AddDate(0, 0, -1)
	case RollWeekly:
		return t.AddDate(0, 0, -7)
	default:
		return t.AddDate(0, 0, -1) // shouldn't occur
	}
}
