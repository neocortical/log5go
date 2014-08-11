log4go
======

A simple, powerful logging library for Go.

Example:
========

```go
import "github.com/neocortical/log4go"

// creates a file logger in the working directory, logging INFO and above
log, err := log4go.NewFileLogger("", "foo.log", log4go.LogInfo, log4go.TF_GoStd)

log.Info("The running progam is called %s", os.Args[0])
log.Debug("This message won't show up because the log level is too low")

// creates a console logger on stdout with no timestamp (you can use any valid Go time format)
log = log4go.NewConsoleLogger(log4go.LogAll, "")

// creates a console logger that logs TRACE,DEBUG,INFO to stdout and WARN,ERROR,FATAL to stderr
log = log4go.NewConsoleLoggerWithStderr(log4go.LogAll, log4go.TF_GoStd)

log.Info("Our PID is %d", os.Getpid())
log.Error("This prints to stderr")

```
