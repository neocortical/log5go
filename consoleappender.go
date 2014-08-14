package log5go

import (
	"io"
	"time"
)

type consoleAppender struct {
	dest io.Writer
	errDest io.Writer
}

func (a *consoleAppender) Append(msg string, level LogLevel, tstamp time.Time) {
	if a.errDest != nil && level >= LogWarn {
		a.errDest.Write([]byte(msg))
	} else {
		a.dest.Write([]byte(msg))
	}
}
