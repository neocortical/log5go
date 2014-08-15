/*
Package log5go is a simple, powerful logging framework for Go, loosely based on Java's log4j.
Loggers support log4j's basic log levels out of the box: TRACE, DEBUG, INFO, WARN, ERROR,
and FATAL, but additional levels can be integrated into the framework while still respecting
the level hierarchy. log5go can log to the console or to one or more files and supports
log file archiving and rotation.

Basics

Loggers are configured using a builder pattern starting with NewLog(LogLevel) and
terminating either with Build() or BuildAndRegister(string). The former simply constructs a
new logger and hands it back to the caller. The latter registers the logger internally
using the caller's source path and supplied key. This results in package-safe loggers
that can be statically retrieved in other parts of the package using GetLog(string).

Examples

The following example creates a file logger and registers it with the name "db":

  log, err := log5go.Log(log5go.LogDebug).ToFile("/var/log", "myprog_db.log").Register("db")

All package local code will be able to retrieve the same logger by calling:

  log, err := log5go.GetLog("db")

This allows logging to be unobtrusive in code, since any package-local code can easily
obtain the desired logger without the need to create a global variable.

The following example creates a file logger with a log rotation scheme:

  log, err := log5go.Log(log5go.LogAll).ToFile("/var/log", "myprog.log").WithRotation(log5go.RollDaily, 7).Build()

In this example, the logger will archive the log file daily at midnight, maintaining a maximum
of 7 archived log files. (A timestamp is appended to the name of each log file and an attempt is
made to delete the file that was created 8 days ago.)
*/
package log5go

// Package version info
const VERSION = "0.8.0"
const MAJOR_VERSION = 0
const MINOR_VERSION = 8
const PATCH_VERSION = 0
