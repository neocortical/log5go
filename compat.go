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

// GoLogger interface implements the Go stdlib log package. Log5Go is almost
// backward-compatible with stdlib. The only differences are A) we implement
// logging ops as an interface instead of a struct so passing around *Logger
// will break, and B) Fatal and Panic ops have been changed to GoFatal and
// GoPanic due to name collisions with the main Log5Go interface (and because
// we don't feel that logging should have the side-effect of terminating a
// program).
type GoLogger interface {
	// Output writes a string to the logger's destination, with the given calldepth applied to identifying the caller
	Output(calldepth int, s string) error
	// SetOutput sets the output destination of the logger. If called on a file logger, a new appender will be created.
	SetOutput(out io.Writer)
	// Flags returns the Go stdlib-specific flags set on the logger.
	Flags() int
	// SetFlags sets the Go stdlib-specifc flags set on the logger.
	SetFlags(flag int)
	// Prefix returns the prefix set on the logger.
	Prefix() string
	// SetPrefix sets the prefix of the logger.
	SetPrefix(prefix string)
	// Print logs a message using the behavior of fmt.Print()
	Print(v ...interface{})
	// Printf logs a message using the behavior of fmt.Printf()
	Printf(format string, v ...interface{})
	// Println logs a message using the behavior of fmt.Println()
	Println(v ...interface{})
	// GoFatal logs a message and calls os.Exit(1)
	GoFatal(v ...interface{})
	// GoFatalf logs a message and calls os.Exit(1)
	GoFatalf(format string, v ...interface{})
	// GoFatalln logs a message and calls os.Exit(1)
	GoFatalln(v ...interface{})
	// GoPanic logs a message and calls panic with the formatted message
	GoPanic(v ...interface{})
	// GoPanicf logs a message and calls panic with the formatted message
	GoPanicf(format string, v ...interface{})
	// GoPanicln logs a message and calls panic with the formatted message
	GoPanicln(v ...interface{})
}

// New creates a new Log5Go with the desired Go stdlib log settings
func New(out io.Writer, prefix string, flag int) Log5Go {
	l := Logger(LogAll).ToWriter(out).WithTimeFmt(parseTimeFmt(flag))
	if prefix != "" {
		l = l.WithPrefix(prefix)
	}
	lines := parseLines(flag)
	if lines == Lshortfile {
		l = l.WithLn()
	} else if lines == Llongfile {
		l = l.WithLine()
	}

	// HACKY: need to set flags directly on the underlying struct
	logger, _ := l.(*logger)
	logger.flag = flag

	return l
}

func (l *logger) Output(calldepth int, s string) error {
	return l.log(time.Now(), LogInfo, calldepth + 1, s, nil)
}

func (l *logger) SetOutput(out io.Writer) {
	l.lock.Lock()
	a, ok := l.appender.(*writerAppender)
	if !ok {
		l.appender = &writerAppender{dest: out, errDest: nil}
	} else {
		a.lock.Lock()
		a.dest = out
		a.errDest = nil
		a.lock.Unlock()
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

// Flags returns the flags on the default logger
func Flags() int {
	return std.Flags()
}

func (l *logger) SetFlags(flag int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.timeFormat = parseTimeFmt(flag)
	l.lines = parseLines(flag)
	l.flag = flag
}

// SetFlags sets the flags on the default logger
func SetFlags(flag int) {
	std.SetFlags(flag)
}

func (l *logger) Prefix() string {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.prefix
}

// Prefix returns the prefix of the default logger
func Prefix() string {
	return std.Prefix()
}

func (l *logger) SetPrefix(prefix string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.prefix = prefix
}

// SetPrefix sets the prefix of the default logger
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.log(time.Now(), LogInfo, 2, fmt.Sprintf(format, v...), nil)
}

func (l *logger) Print(v ...interface{}) {
	l.log(time.Now(), LogInfo, 2, fmt.Sprint(v...), nil)
}

func (l *logger) Println(v ...interface{}) {
	l.log(time.Now(), LogInfo, 2, fmt.Sprintln(v...), nil)
}

func (l *logger) GoFatal(v ...interface{}) {
	l.log(time.Now(), LogFatal, 2, fmt.Sprint(v...), nil)
	os.Exit(1)
}

func (l *logger) GoFatalf(format string, v ...interface{}) {
	l.log(time.Now(), LogFatal, 2, fmt.Sprintf(format, v...), nil)
	os.Exit(1)
}

func (l *logger) GoFatalln(v ...interface{}) {
	l.log(time.Now(), LogFatal, 2, fmt.Sprintln(v...), nil)
	os.Exit(1)
}

func (l *logger) GoPanic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.log(time.Now(), LogFatal, 2, s, nil)
	panic(s)
}

func (l *logger) GoPanicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.log(time.Now(), LogFatal, 2, s, nil)
	panic(s)
}

func (l *logger) GoPanicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.log(time.Now(), LogFatal, 2, s, nil)
	panic(s)
}

// Print calls Print on the default logger
func Print(v ...interface{}) {
	std.log(time.Now(), LogInfo, 2, fmt.Sprint(v...), nil)
}

// Printf calls Printf on the default logger
func Printf(format string, v ...interface{}) {
	std.log(time.Now(), LogInfo, 2, fmt.Sprintf(format, v...), nil)
}

// Println calls Println on the default logger
func Println(v ...interface{}) {
	std.log(time.Now(), LogInfo, 2, fmt.Sprintln(v...), nil)
}

// GoFatal calls GoFatal on the default logger
func GoFatal(v ...interface{}) {
	std.log(time.Now(), LogFatal, 2, fmt.Sprint(v...), nil)
	os.Exit(1)
}

// GoFatalf calls GoFatalf on the default logger
func GoFatalf(format string, v ...interface{}) {
	std.log(time.Now(), LogFatal, 2, fmt.Sprintf(format, v...), nil)
	os.Exit(1)
}

// GoFatalln calls GoFatalln on the default logger
func GoFatalln(v ...interface{}) {
	std.log(time.Now(), LogFatal, 2, fmt.Sprintln(v...), nil)
	os.Exit(1)
}

// GoPanic calls GoPanic on the default logger
func GoPanic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.log(time.Now(), LogFatal, 2, s, nil)
	panic(s)
}

// GoPanicf calls GoPanicf on the default logger
func GoPanicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.log(time.Now(), LogFatal, 2, s, nil)
	panic(s)
}

// GoPanicln calls GoPanicln on the default logger
func GoPanicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.log(time.Now(), LogFatal, 2, s, nil)
	panic(s)
}

// parseTimeFmt extracts a time format string from Go log flags
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

// parseLines extracts an isolated Go std log flag (can be 0, Lshortfile, or Llongfile)
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
