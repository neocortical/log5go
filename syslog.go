package log5go

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var socketTypes = []string{"unixgram", "unix"}

var socketLocations = []string{"/dev/log", "/var/run/syslog", "/var/run/log"}

// SyslogPriority represents a syslog priority level
type SyslogPriority int

// syslog severity levels
const (
	SyslogEmergency SyslogPriority = iota
	SyslogAlert
	SyslogCritical
	SyslogError
	SyslogWarning
	SyslogNotice
	SyslogInfo
	SyslogDebug
)

// syslog facility levels
const (
	SyslogKernel SyslogPriority = iota << 3
	SyslogUser
	SyslogMail
	SyslogDaemon
	SyslogAuth
	SyslogSyslog
	SyslogLpr
	SyslogNews
	SyslogUUCP
	SyslogClock
	SyslogAuthpriv
	SyslogFTP
	SyslogNTP
	SyslogLogAudit
	SyslogLogAlert
	SyslogCron
	SyslogLocal0
	SyslogLocal1
	SyslogLocal2
	SyslogLocal3
	SyslogLocal4
	SyslogLocal5
	SyslogLocal6
	SyslogLocal7
)

type syslogAppender struct {
	sync.Mutex
	conn     net.Conn
	facility SyslogPriority
	tag      string
	hostname string
}

func (a *syslogAppender) Append(msg *[]byte, level LogLevel, tstamp time.Time) error {
	a.Lock()
	defer a.Unlock()

	TerminateMessageWithNewline(msg)

	// sanity check
	if a.conn == nil {
		return fmt.Errorf("connection is not established")
	}

	pri := calculatePriority(a.facility, level)
	hostname := a.calculateHostname()

	_, err := fmt.Fprintf(a.conn, "<%d>1 %s %s %s[%d]: %s", pri, tstamp.Format(time.RFC3339Nano), hostname, a.tag, os.Getpid(), *msg)
	return err
}

func (a *syslogAppender) calculateHostname() string {
	if a.hostname != "" {
		return a.hostname
	}

	// try to get local IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		// unlikely to happen, but nothing we can do
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				a.hostname = ipnet.IP.String()
				return a.hostname
			}
		}
	}

	// try to use hostname if we can't get an IP addr
	host, err := os.Hostname()
	if err == nil {
		a.hostname = host
		return a.hostname
	}

	// everything failed. just put something in a.addr so we don't call method over and over
	a.hostname = "unknown-host"
	return a.hostname
}

type syslogFormatter struct {
	formatter *StringFormatter
}

func newSyslogFormatter(lines bool) Formatter {
	var inner *StringFormatter
	if lines {
		inner = NewStringFormatter("%p (%c:%n): %m")
	} else {
		inner = NewStringFormatter("%p: %m")
	}

	return &syslogFormatter{formatter: inner}
}

func (f *syslogFormatter) Format(tstamp time.Time, level LogLevel, prefix, caller string, line uint, msg string, data Data, out *[]byte) {
	f.formatter.Format(tstamp, level, prefix, caller, line, msg, data, out)
}

func (f *syslogFormatter) SetTimeFormat(timeFormat string) {
	// NOOP
}

func (f *syslogFormatter) SetLines(lines bool) {
	f.formatter.SetLines(lines)
}

func calculatePriority(facility SyslogPriority, level LogLevel) SyslogPriority {
	switch {
	case level <= LogDebug:
		return facility | SyslogDebug
	case level <= LogInfo:
		return facility | SyslogInfo
	case level <= LogNotice:
		return facility | SyslogNotice
	case level <= LogWarn:
		return facility | SyslogWarning
	case level <= LogError:
		return facility | SyslogError
	case level <= LogCritical:
		return facility | SyslogCritical
	case level <= LogAlert:
		return facility | SyslogAlert
	default:
		return facility | SyslogEmergency

	}
}
