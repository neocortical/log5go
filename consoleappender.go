package log5go

import (
	"os"
	"time"
)

type consoleAppender struct {
	stderrAware bool
}

func (a *consoleAppender) Append(msg string, level LogLevel, tstamp time.Time) {
	if a.stderrAware && level >= LogWarn {
		os.Stderr.Write([]byte(msg))
	} else {
		os.Stdout.Write([]byte(msg))
	}
}
