package log5go

type Formatter interface {
	Format(timeString, levelString, prefix, caller string, line uint, msg string, data Data) []byte
}

// Some constant string formats for convenience, also used internally
const (
	FMT_Default = "%t %l : %m"
	FMT_DefaultPrefix = "%t %l %p: %m"
	FMT_DefaultLines = "%t %l (%c:%n): %m"
	FMT_DefaultPrefixLines = "%t %l %p (%c:%n): %m"
	FMT_NoTime							= "%l : %m"
	FMT_NoTimePrefix = "%l %p: %m"
	FMT_NoTimeLines = "%l (%c:%n): %m"
	FMT_NoTimePrefixLines = "%l %p (%c:%n): %m"
)

// defaults, use dependent on logger settings
var fmtNone = NewStringFormatter(FMT_NoTime)
var fmtTime = NewStringFormatter(FMT_Default)
var fmtTimePrefix = NewStringFormatter(FMT_DefaultPrefix)
var fmtTimeLines = NewStringFormatter(FMT_DefaultLines)
var fmtTimePrefixLines = NewStringFormatter(FMT_DefaultPrefixLines)

var fmtNotimePrefix = NewStringFormatter(FMT_NoTimePrefix)
var fmtNotimeLines = NewStringFormatter(FMT_NoTimeLines)
var fmtNotimePrefixLines = NewStringFormatter(FMT_NoTimePrefixLines)
