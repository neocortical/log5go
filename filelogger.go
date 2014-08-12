package log4go

type RollFrequency uint8

// Log rotation frequencies. Daily rotates at midnight, weekly rotates on Sunday at midnight
const (
	RollNone     RollFrequency = iota
	RollMinutely               // nice for testing
	RollHourly
	RollDaily
	RollWeekly
)

const SaveAllOldLogs = -1
