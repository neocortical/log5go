package log5go

import (
	"io"
	"sync"
	"time"
)

type writerAppender struct {
	lock sync.Mutex
	dest io.Writer
	errDest io.Writer
}

func (a *writerAppender) Append(msg string, level LogLevel, tstamp time.Time) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.errDest != nil && level >= LogWarn {
		a.errDest.Write([]byte(msg))
	} else {
		a.dest.Write([]byte(msg))
	}
}
