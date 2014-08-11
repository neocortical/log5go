log4go
======

A simple logging library for Go.

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
log = log4go.NewConsoleLogger(log4go.LogAll, log4go.TF_GoStd)
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A console logger that writes errors to stderr
---------------------------------------------

```go
log = log4go.NewConsoleLoggerWithStderr(log4go.LogAll, log4go.TF_GoStd)
log.Info("Trace, debug, and info go to stdout")
log.Error("Warn, error, and fatal go to stderr")
```

A simple file logger
--------------------

```go
log, err := log4go.NewFileLogger("/tmp", "foo.log", log4go.LogInfo, log4go.TF_GoStd)
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A rolling file appender
-----------------------

```go
log, err := log4go.NewRollingFileLogger("/tmp", "foo.log", log4go.LogInfo, log4go.TF_GoStd, log4go.RollDaily, 10)
log.Info("Hello, World. My PID is %d", os.Getpid())
```

Custom logging levels
---------------------
```go
var LLCustomDebug log4go.LogLevel = log4go.LogDebug + 1
var LLCustomLogLevel log4go.LogLevel = log4go.LogInfo + 1
var LLCustomInfo log4go.LogLevel = log4go.LogInfo + 2

log := log4go.NewConsoleLogger(LLCustomLogLevel, log4go.TF_GoStd)

log.Log(LLCustomDebug, "Won't see this")
log.Info("Won't see this either")
log.Log(LLCustomInfo, "This will get logged")
```

Features
========

* Dead simple: Create a logger and go
* Supports string formatting, just like fmt.Printf()
* Standard built-in log levels: TRACE, DEBUG, INFO, WARN, ERROR, FATAL
* Console or file logging
* Interleave custom log levels with standard ones 
* Full control over date/time format (uses time.Format under the hood)
* Rolling file appender (roll each minute, hour, day, or week)
* Console log can send errors to stderr instead of stdout

TODO
====

* Testing!
* Scheduled log rotation (currently, logs only roll when a message arrives)
* Delete old logs (currently, the oldLogsToSave value is ignored)

Caveats
=======

Log4Go is something I whipped up in a day, because I was unsatisfied with the Go default logging library. I do not recommend using Log4Go for production code until it's been properly tested.
