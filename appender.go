package log5go

import (
	"time"
)

// Appender interface is responsible for writing a formatted message to
// a destination. Log5Go provides rich support for appending to files, the console,
// and to arbitrary writers. Developers can extend log5go by implementing a custom
// Appender and configuring their logger with ToAppender(a).
type Appender interface {
	Append(msg *[]byte, level LogLevel, tstamp time.Time) error
}

// TerminateMessageWithNewline function tests msg content and adds a terminating
// newline if not there already. If you write a custom appender and want line
// termination, you should call this function on the msg before writing it.
func TerminateMessageWithNewline(msg *[]byte) {
	if len(*msg) == 0 || (*msg)[len(*msg)-1] != '\n' {
		*msg = append(*msg, '\n')
	}
}
