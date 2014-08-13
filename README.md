log5go
======

(Yet another) simple, powerful logging library for Go.

Very loosely based on the (in)famous log4j Java logging library, but uncluttered, Go-like, and awesome.

Why 5? Because there are already about a squillion log4go projects on Github and I wanted to take it to the next level!


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

TODO
====

* More testing. There's coverage for time (including DST rotation issues), log builder, and regsistry, but need more.
* Custom layouts (pattern, JSON, HTML, etc.)
* syslog
* custom (third-party) appenders

Caveats
=======

log5go is a young project and has not been deployed in production environments. Use at your own risk.

Please feel free to contribute feedback, advice, feature suggestions, and pull requests!
