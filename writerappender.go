package log5go

import (
	"io"
	"sync"
	"time"
)

type writerAppender struct {
	lock    sync.Mutex
	dest    io.Writer
	errDest io.Writer
}

func (a *writerAppender) Append(msg *[]byte, level LogLevel, tstamp time.Time) (err error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	TerminateMessageWithNewline(msg)

	if a.errDest != nil && level >= LogWarn {
		_, err = a.errDest.Write(*msg)
	} else {
		_, err = a.dest.Write(*msg)
	}

	return err
}
