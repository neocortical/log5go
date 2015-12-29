log5go
======

(Yet another) simple, powerful logging library for Go.

Very loosely based on the (in)famous log4j Java logging library, but uncluttered, Go-like, and awesome.

Log5Go is no longer compatible with Go's log package. I believe that Log5Go is an improvement on Go's logging paradigm for the following reasons:
* Additional functionality such as registry, log levels, data binding, and log file rotation
* The Go logger's Fatal and Panic methods have the side effect of stopping the program, which I believe is inappropriate for a logging package


[![Build Status](https://travis-ci.org/neocortical/log5go.svg?branch=master)](https://travis-ci.org/neocortical/log5go) [![Coverage](http://gocover.io/_badge/github.com/neocortical/log5go?v=1)](http://gocover.io/github.com/neocortical/log5go?v=1) [![GoDoc](https://godoc.org/github.com/neocortical/log5go?status.svg)](https://godoc.org/github.com/neocortical/log5go)

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

A simple console logger (defaults to stderr)
--------------------------------------------

```go

log := l5g.Logger(l5g.LogAll)
log.Info("Hello, World. My PID is %d", os.Getpid())

```

A logger to stdout with a custom time format
--------------------------------------------

```go

log = l5g.Logger(l5g.LogDebug).ToStdout().WithTimeFmt("Jan _2 15:04:05")
log.Info("Hello, World. My PID is %d", os.Getpid())

```

A console logger that writes to stdout but sends errors to stderr
-----------------------------------------------------------------

```go

log = l5g.Logger(l5g.LogAll).ToStdout().WithStderr()
log.Info("Trace, debug, and info go to stdout")
log.Error("Warn, error, and fatal go to stderr")

```

A JSON logger
-------------

```go

log = l5g.Logger(l5g.LogAll).ToStdout().Json()
log.Info("I'm inside a JSON string!")

```

A logger with structured data
-----------------------------

```go

log = l5g.Logger(l5g.LogAll).ToStdout()
log.WithData(l5g.Data{"foo":"bar", "baz":1}).Info("Hey look, some data: ")

```

A simple file logger
--------------------

```go

log = l5g.Logger(l5g.LogInfo).ToFile("/tmp", "foo.log")
log.Info("Hello, World. My PID is %d", os.Getpid())

```

A rolling file appender
-----------------------

```go

log = l5g.Logger(l5g.LogDebug).ToFile("/tmp", "foo.log").WithRotation(l5g.RollDaily, 7)
log.Info("Hello, World. My PID is %d", os.Getpid())

```

Custom log output format
------------------------

```go

log = l5g.Logger(l5g.LogDebug).WithFmt("%m") // message only
log.Info("He hates these cans!")

```

Using the internal registry
---------------------------

```go

// In mypkg/foo.go
log, err := l5g.Logger(l5g.LogDebug).ToFile("/tmp", "mypkg.log").Register("mypkg/mainlog")
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

log = l5g.Logger(LLCustomLogLevel)

log.Log(LLCustomDebug, "Won't see this: priority too low")
log.Info("Won't see this either")
log.Log(LLCustomInfo, "This will get logged with the prefix we registered")

```

Syslog
------

```go

// local (Unix socket) syslogd
log = l5g.Logger(LogDebug).ToLocalSyslog(l5g.SyslogLocal2, "myapp")

// remote syslog (supports "tcp", "udp", Unix sockets)
log = l5g.Logger(LogDebug).ToRemoteSyslog(l5g.SyslogLocal2, "myapp", "tcp", "syslogd.example.com:514")

```
Note: All syslog logging priorities are supported. Fatal() will log as EMERG. Otherwise, the naming is 1:1.

Default Logger
--------------

A default logger is availabe via `GetLogger(key string)` and is configured from environment variables. If no log5go environment variables are set this logger defaults to logging to stdout and stderr. Log messages are prefixed with `key`.

Example:

```go
log = l5g.GetLogger("foo")

log.Info("some information worth logging.")
```

The following environment variables can be used to configure default loggers:
* L5G_LOG_FILE_NAME - The full path to the logfile to log to. If this environment variable is not set the default logger will log to stdout and stderr.
* L5G_LOG_LEVEL - The level: ALL, TRACE, DEBUG, INFO, NOTICE, WARN, ERROR, CRIT, ALERT, or FATAL.
* L5G_LOG_LINE_LENGTH - LONG, SHORT, or NONE. SHORT includes the name of the source file and line number, LONG includes the full path to source and line number, NONE excludes source and line number information in log messages.
* L5G_LOG_FILE_ROTATION_FREQUENCY - NONE, MINUTE, HOUR, DAY or WEEK. Frequency to rotate log files.
* L5G_LOG_FILE_ROTATION_KEEP_N_FILES - Number of previous log files to keep.


Features
========

* Dead simple: Build a logger and go
* Supports string formatting, just like fmt.Printf()
* Standard built-in log levels: TRACE, DEBUG, INFO, WARN, ERROR, FATAL
* Console or file logging
* JSON layouts for consumption by scripts, Splunk, etc.
* Add custom structured state data to log messages
* Syslog support
* Register loggers and retrieve by key (no globals, or passing logs around)
* Interleave custom log levels with standard ones
* Full control over date/time format (uses time.Format under the hood)
* Rolling file appender (roll each minute, hour, day, or week)
* Optionally store N old log files with date stamps
* Console log can send errors to stderr instead of stdout
* Extensible through custom appenders


Log Format
==========

The default output format for log messages consists the following
ordering of elements: timestamp, level, prefix, caller:line, message. Depending on how
you set up your logger, not all of these elements will be present. Here are some examples of
log output:

```

# Default log format:
2009/01/23 01:23:23 INFO : I'm picking out a Thermos for you...

# Everything:
2009/01/23 01:23:23.123123 INFO myprefix (acme.go:123): Not an ordinary Thermos, for you...

# No time or line info:
2009/01/23 INFO myprefix: But the extra best thermos you can buy...

# No prefix:
2009/01/23 01:23:23 INFO (acme.go:123): With vinyl and stripes...

# No level string for level:
2009/01/23 01:23:23 : And a cup built right in.

# No nothing (.WithFmt("%m")):
He hates these cans! Stay away from the cans!

```

You can also log in JSON format by calling the .Json() method on a logger.

Go's stdlib log functions print at level INFO and GoPanic() and GoFatal() print at level FATAL.


Performance
===========

I've done some basic performance tuning, with the result that simple logging is
reasonably close to Go's pkg/log logger. Here are the current performance results:

### Test setup:

```go
l5g.Logger(l5g.LogAll).WithPrefix("l5g").ToFile("/tmp", "l5gtest.log")
// vs
f, _ := os.Create("/tmp/gologtest.log")
log.New(f, "stdlib", log.LstdFlags)
```

* Single goroutine logs 10M lines in a loop

### Environment:

```
MBP, 2.3 GHz Intel i7, 8GB 1333MHz DDR3, OS X 10.9.4
```

### Results:

```
pkg/log: real: 33.24, user: 12.03, sys: 21.12, Max heap: 591KB, Total bytes allocated: 160MB, GCs: 570

log5go:  real: 47.40, user: 22.82, sys: 24.46, Max heap: 582KB, Total bytes allocated: 480MB, GCs: 1611
```

### Conclusion:

Log5Go is reasonably fast and memory-efficient and in the same ballpark as pkg/log but
allocates much more memory and could be improved further.


ROADMAP
=======

* HTML structured layout
* Futher performance benchmarking and tuning
* log chaining?
* 1.0 (stable API) release


About the Developer
===================

Nathan Smith (neocortical) is a seasoned Java architect and developer and a lover of Go.
Email him at nathan@neocortical.net

Please feel free to contribute feedback, advice, feature suggestions, and pull requests!
