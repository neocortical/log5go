log5go
======

(Yet another) simple, powerful logging library for Go.

Very loosely based on the (in)famous log4j Java logging library, but uncluttered, Go-like, and awesome.

Log5Go is *mostly* compatible with Go's log package, and is almost a drop-in replacement. See section on replacing the Go logger for details.

[![Build Status](https://travis-ci.org/neocortical/log5go.svg?branch=master)](https://travis-ci.org/neocortical/log5go)

Install
=======

```
go get github.com/neocortical/log5go
```

And import:
```go
import l5g "github.com/neocortical/log5go"
```

Examples
========

A simple console logger
-----------------------

```go
log, err := l5g.Log(l5g.LogAll).ToStdout().Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A logger with a custom time format
----------------------------------

```go
log, err = l5g.Log(l5g.LogDebug).ToStdout().WithTimeFmt("Jan _2 15:04:05").Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A console logger that writes errors to stderr
---------------------------------------------

```go
log, err = l5g.Log(l5g.LogAll).ToStdout().WithStderr().Build()
log.Info("Trace, debug, and info go to stdout")
log.Error("Warn, error, and fatal go to stderr")
```

A simple file logger
--------------------

```go
log, err = l5g.Log(l5g.LogInfo).ToFile("/tmp", "foo.log").Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

A rolling file appender
-----------------------

```go
log, err = l5g.Log(l5g.LogDebug).ToFile("/tmp", "foo.log").WithRotation(l5g.RollDaily, 7).Build()
log.Info("Hello, World. My PID is %d", os.Getpid())
```

Using the internal registry
---------------------------

```go
// In mypkg/foo.go
log, err := l5g.Log(l5g.LogDebug).ToFile("/tmp", "mypkg.log").Register("mypkg/mainlog")
log.Info("Hello from file foo.go")

// In mypkg/bar.go
log, err := l5g.GetLog("mypkg/mainlog")
log.Info("Hello from file bar.go")
```
Note: It's a good convention to prefix log names with your package name to avoid collisions when
more than one package uses log5go in the same process.

Custom logging levels
---------------------
```go
var LLCustomDebug l5g.LogLevel = l5g.LogDebug + 1
var LLCustomLogLevel l5g.LogLevel = l5g.LogInfo + 1
var LLCustomInfo l5g.LogLevel = l5g.LogInfo + 2

l5g.RegisterLogLevel(LLCustomInfo, "CUSTOM_INFO") // optional

log, err = l5g.Log(LLCustomLogLevel).ToStdout().Build()

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
* Extensible through custom appenders

Replacing Go's Logger
=====================

Log5Go is almost a drop-in replacement for Go's built-in log package. Due to some of
Log5Go's design decisions, there are a couple of incompatibilities that may need to be
dealt with. Here are the steps to upgrade from the standard log package to Log5Go:

* Change all imports from import "log" to import log "github.com/neocortical/log5go"
* Change all calls to Panic(), Fatal(), and their variants to GoPanic(), GoFatal(), etc. (Print() and similar are fine)
* If your code passes logs around, change all references from *Logger to Log5Go


TODO
====

* More testing. There's coverage for time (including DST rotation issues), log builder, and regsistry, but need more.
* Custom layouts (pattern, JSON, HTML, etc.)
* syslog support (better than Go's native support, about which even Go authors just shake their heads)
* log chaining

Caveats
=======

log5go is a young project and has not been deployed in production environments. Use at your own risk.

Please feel free to contribute feedback, advice, feature suggestions, and pull requests!
