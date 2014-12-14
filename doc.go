/*
Package log5go is a powerful, configurable logging framework for Go.

log5go features clear, declarative configuration, 6 standard log levels (based on Java's Log4j),
logfile rotation, and highly configurable formats and appenders (file, console, Writer, etc.).

log5go loggers can be added to and retrieved from a registry, removing the need to
pass loggers around your code or rely on global logger variables.


Basics

Loggers are configured using a builder pattern starting with Logger(LogLevel) followed by
any of a rich set of configuration methods. Want JSON? Use .Json(). Want logging to a file
with daily rotation? Use .ToFile(dir, fname).WithRotation(frequency, numSaved)

You can configure a number of different format and destination options with a single,
easy-to-read line of code.


Examples

Create a logger that logs JSON to stderr:

  log, err := log5go.Logger(log5go.LogInfo).Json()

Create a file logger with an hourly log rotation and 10 saved archive files:

  log, err := log5go.Logger(log5go.LogDebug).ToFile("/var/log", "db.log").WithRotation(log5go.RollHourly, 10)

Register a stdout logger and retrieve it in another part of your code:

  // create and register
  log, err := log5go.Logger(log5go.LogAll).ToStdout().Register("mylog")

  // get it from the registry
  log, err := log5go.GetLog("mylog")

*/
package log5go

// Package version info
const VERSION = "0.14.0"
const MAJOR_VERSION = 0
const MINOR_VERSION = 15
const PATCH_VERSION = 0
