log4go
======

A simple, powerful logging library for Go.

Install
=======

```
go get github.com/neocortical/log4go 
```

And import: 
```go
import "github.com/neocortical/log4go"
```

Examples
========

A simple console logger
-----------------------

```go
log, err := log4go.NewLog(log4go.LogAll).ToConsole().Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A logger with a custom time format
----------------------------------

```go
log, err = log4go.NewLog(log4go.LogAll).ToConsole().WithTimeFormat("Jan _2 15:04:05").Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A console logger that writes errors to stderr
---------------------------------------------

```go
log, err = log4go.NewLog(log4go.LogAll).ToConsole().WithStderrSupport().Build()
log.Info("Trace, debug, and info go to stdout")
log.Error("Warn, error, and fatal go to stderr")
```

A simple file logger
--------------------

```go
log, err = log4go.NewLog(log4go.LogInfo).ToFile("/tmp", "foo.log").Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A rolling file appender
-----------------------

```go
log, err = log4go.NewLog(log4go.LogDebug).ToFile("/tmp", "foo.log").WithFileRotation(log4go.RollDaily, 7).Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

Custom logging levels
---------------------
```go
var LLCustomDebug log4go.LogLevel = log4go.LogDebug + 1
var LLCustomLogLevel log4go.LogLevel = log4go.LogInfo + 1
var LLCustomInfo log4go.LogLevel = log4go.LogInfo + 2

log4go.RegisterLogLevel(LLCustomInfo, "CUSTOM_INFO") // optional

log, err = log4go.NewLog(LLCustomLogLevel).ToConsole().Build()

log.Log(LLCustomDebug, "Won't see this: priority too low")
log.Info("Won't see this either")
log.Log(LLCustomInfo, "This will get logged with the prefix we registered")
```

Features
========

* Dead simple: Build a logger and go
* Supports string formatting, just like fmt.Printf()
* Standard built-in log levels: TRACE, DEBUG, INFO, WARN, ERROR, FATAL
* Console or file logging
* Interleave custom log levels with standard ones 
* Full control over date/time format (uses time.Format under the hood)
* Rolling file appender (roll each minute, hour, day, or week)
* Optionally store N old log files with date stamps
* Console log can send errors to stderr instead of stdout

TODO
====

* Testing! Pretty good coverage for log rotation date math, but litte else
* Log registry for "static" retrieval of loggers via registered names
* Custom layouts

Caveats
=======

Log4Go is something I whipped up in a couple days, because I was unsatisfied with the Go default logging library. I do not recommend using Log4Go for production code until it's been properly vetted.

Please feel free to contribute feedback, advice, and pull requests! 
