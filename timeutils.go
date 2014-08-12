package log4go

import (
  "time"
)

// given a time, calculate the instant that the log should next roll
func calculateNextRollTime(t time.Time, freq RollFrequency) time.Time {
  if freq == RollMinutely {
    t = t.Truncate(time.Minute)
    return t.Add(time.Minute)
  } else if freq == RollHourly {
    t = t.Truncate(time.Hour)
    return t.Add(time.Hour)
  } else {
    t = t.Truncate(time.Hour)
    t = t.Add(-time.Hour*time.Duration(t.Hour()))
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
func calculatePreviousRollTime(t time.Time, freq RollFrequency) time.Time {
  switch freq {
  case RollMinutely:
    return t.Add(-time.Minute)
  case RollHourly:
    return t.Add(-time.Hour)
  case RollDaily:
    return t.AddDate(0, 0, -1)
  case RollWeekly:
    return t.AddDate(0, 0, -7)
  default:
    return t.AddDate(0, 0, -1) // shouldn't occur
  }
}
