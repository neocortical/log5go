log4go
======

A simple, powerful logging library for Go.

Very loosely based on the (in)famous log4j Java logging library.

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
log, err := log4go.Log(log4go.LogAll).ToConsole().Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A logger with a custom time format
----------------------------------

```go
log, err = log4go.Log(log4go.LogDebug).ToConsole().WithTimeFmt("Jan _2 15:04:05").Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A console logger that writes errors to stderr
---------------------------------------------

```go
log, err = log4go.Log(log4go.LogAll).ToConsole().WithStderr().Build()
log.Info("Trace, debug, and info go to stdout")
log.Error("Warn, error, and fatal go to stderr")
```

A simple file logger
--------------------

```go
log, err = log4go.Log(log4go.LogInfo).ToFile("/tmp", "foo.log").Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A rolling file appender
-----------------------

```go
log, err = log4go.Log(log4go.LogDebug).ToFile("/tmp", "foo.log").WithRotation(log4go.RollDaily, 7).Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

Using the internal registry
---------------------------

```go
// In mypkg/foo.go
log, err := log4go.Log(log4go.LogDebug).ToFile("/tmp", "mypkg.log").Register("mypkg/mainlog")
log.Info("Hello from file foo.go")

// In mypkg/bar.go
log, err := log4go.GetLog("mypkg/mainlog")
log.Info("Hello from file bar.go")
```
Note: It's a good convention to prefix log names with your package name to avoid collisions when
more than one package uses Log4Go in the same process.

Custom logging levels
---------------------
```go
var LLCustomDebug log4go.LogLevel = log4go.LogDebug + 1
var LLCustomLogLevel log4go.LogLevel = log4go.LogInfo + 1
var LLCustomInfo log4go.LogLevel = log4go.LogInfo + 2

log4go.RegisterLogLevel(LLCustomInfo, "CUSTOM_INFO") // optional

log, err = log4go.Log(LLCustomLogLevel).ToConsole().Build()

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
* Register loggers and retrieve by key (no globals, or passing logs around)
* Interleave custom log levels with standard ones
* Full control over date/time format (uses time.Format under the hood)
* Rolling file appender (roll each minute, hour, day, or week)
* Optionally store N old log files with date stamps
* Console log can send errors to stderr instead of stdout

TODO
====

* Testing! Pretty good coverage for log rotation date math, but needs builder, appender test coverage
* Custom layouts (pattern, JSON, HTML, etc.)

Caveats
=======

Log4Go is a young project and has not been deployed in production environments. Use at your own risk.

Please feel free to contribute feedback, advice, feature suggestions, and pull requests!
