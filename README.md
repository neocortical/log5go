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
* Log output format is different from Go's. See the next section for details


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

Go's stdlib log functions print at level INFO and GoPanic() and GoFatal() print at level FATAL.


TODO
====

* Structured layouts (JSON, HTML, etc.)
* syslog support (better than Go's native support)
* log chaining


Caveats
=======

log5go is a young project and has not been deployed in production environments. Use at your own risk.

Please feel free to contribute feedback, advice, feature suggestions, and pull requests!
