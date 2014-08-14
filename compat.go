// This file contains code taken from the Go source code.
// Copyright from original file included here:
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the GO_LICENSE file.

package log5go

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

// Prefexes from stdlib log.go. See Go source for implementation details
const (
	Ldate         = 1 << iota     // the date: 2009/01/23
	Ltime                         // the time: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type GoLogger interface {
	Output(calldepth int, s string) error
	SetOutput(out io.Writer)
	Flags() int
	SetFlags(flag int)
	Prefix() string
	SetPrefix(prefix string)
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	GoFatal(v ...interface{})
	GoFatalf(format string, v ...interface{})
	GoFatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

func New(out io.Writer, prefix string, flag int) Log5Go {
	b := Log(LogAll).ToWriter(out).WithTimeFmt(parseTimeFmt(flag))
	if prefix != "" {
		b = b.WithPrefix(prefix)
	}
	lines := parseLines(flag)
	if lines == Lshortfile {
		b = b.WithLn()
	} else if lines == Llongfile {
		b = b.WithLine()
	}

	l, _ := b.Build()

	// HACKY: need to set flags directly on the underlying struct
	logger, _ := l.(*logger)
	logger.flag = flag

	return l
}

func (l *logger) Output(calldepth int, s string) error {
	now := time.Now() // get this early.
	var file string
	var line int

	l.lock.Lock()
	defer l.lock.Unlock()

	if l.lines != 0 {
		// release lock while getting caller info - it's expensive.
		l.lock.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.lock.Lock()
	}

	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, s...) // TODO: Appender should take []byte
	if len(s) > 0 && s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	// TODO: proper log levels for Output() from GoFatal, Panic
	l.appender.Append(string(l.buf), LogInfo, now)
	return nil // TODO: Appender should return error
}

func (l *logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	*buf = append(*buf, l.prefix...)
	if l.timeFormat != "" {
		*buf = append(*buf, t.Format(l.timeFormat)...)
		*buf = append(*buf, ' ')
	}

	if l.lines&(Lshortfile|Llongfile) != 0 {
		if l.lines == Lshortfile {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, fmt.Sprintf("%s:%d: ", file, line)...)
	}
}

func (l *logger) SetOutput(out io.Writer) {
	l.lock.Lock()
	a, ok := l.appender.(*consoleAppender)
	if !ok {
		l.appender = &consoleAppender{out, nil}
	} else {
		a.dest = out
		a.errDest = nil
	}
	l.lock.Unlock()
}

func SetOutput(out io.Writer) {
	std.SetOutput(out)
}

func (l *logger) Flags() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.flag
}

func (l *logger) SetFlags(flag int) {
	l.lock.Lock()
	l.timeFormat = parseTimeFmt(flag)
	l.lines = parseLines(flag)
	l.flag = flag
	l.lock.Unlock()
}

func SetFlags(flag int) {
	std.SetFlags(flag)
}

func (l *logger) Prefix() string {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.prefix
}

func Prefix() string {
	return std.Prefix()
}

func (l *logger) SetPrefix(prefix string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.prefix = prefix
}

func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *logger) Print(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...))
}

func (l *logger) Println(v ...interface{}) {
	l.Output(2, fmt.Sprintln(v...))
}

func (l *logger) GoFatal(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *logger) GoFatalf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *logger) GoFatalln(v ...interface{}) {
	l.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func (l *logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(2, s)
	panic(s)
}

func (l *logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, s)
	panic(s)
}

func (l *logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(2, s)
	panic(s)
}

func Print(v ...interface{}) {
	std.Output(2, fmt.Sprint(v...))
}

func Printf(format string, v ...interface{}) {
	std.Output(2, fmt.Sprintf(format, v...))
}

func Println(v ...interface{}) {
	std.Output(2, fmt.Sprintln(v...))
}

func Fatal(v ...interface{}) {
	std.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	std.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	std.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(2, s)
	panic(s)
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Output(2, s)
	panic(s)
}

func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Output(2, s)
	panic(s)
}

func parseTimeFmt(flag int) string {
	if flag & (Ldate | Ltime | Lmicroseconds) == 0 {
		return ""
	}

	buf := new(bytes.Buffer)
	if flag & Ldate > 0 {
		buf.Write([]byte("2006/01/02"))
	}
	if flag & (Ltime | Lmicroseconds) > 0 {
		if buf.Len() > 0 {
			buf.Write([]byte(" 15:04:05"))
		} else {
			buf.Write([]byte("15:04:05"))
		}

		if flag & Lmicroseconds > 0 {
			buf.Write([]byte(".000000"))
		}
	}

	return buf.String()
}

func parseLines(flag int) int {
	if flag & (Llongfile | Lshortfile) > 0 {
		if flag & Lshortfile == 0 {
			return Llongfile
		} else {
			return Lshortfile
		}
	}
	return 0
}
