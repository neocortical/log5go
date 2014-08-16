package log5go

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// StringFormatter formats log messages according to a string pattern. The
// pattern consists of literal text, augmented by the following meta-patterns:
//
// %t - time string (as formatted by timeFormat)
// %l - log level (as discovered in registered level strings)
// %p - user-supplied prefix
// %c - caller (if line info present)
// %n - line number (if line info present)
// %m - caller-supplied log message
// %% - literal percent sign
//
// Single occurrences of % will be discarded. Be sure to include %m somewhere or
// your message won't get logged!
type StringFormatter struct {
	parts []string
}

func NewStringFormatter(pattern string) (sf *StringFormatter) {
	sf = new(StringFormatter)
	var buf []byte
	r := make([]byte, 4)
	for len(pattern) > 0 {
		runeValue, width := utf8.DecodeRuneInString(pattern)
		pattern = pattern[width:]
		if runeValue == utf8.RuneError {
			continue
		}

		// intercept meta-pattern. ignore % if meta-pattern is illegal
		if runeValue == '%' {

			// collect the meta-pattern
			meta, width := utf8.DecodeRuneInString(pattern)
			switch meta {
			case '%', 't', 'l', 'p', 'c', 'n', 'm':
			  // valid meta-pattern detected. dump any collected literal pattern first
				// dump any literal value we have collected
				if len(buf) > 0 {
					sf.parts = append(sf.parts, string(buf))
					buf = buf[:0]
				}

				sf.parts = append(sf.parts, "%" + string(meta & 0xff)) // all metas are ascii
				pattern = pattern[width:]
			}
		} else {
			utf8.EncodeRune(r, runeValue)
			buf = append(buf, r[0:width]...)
		}
	}

	if len(buf) > 0 {
		sf.parts = append(sf.parts, string(buf))
	}

	return sf
}

func (f *StringFormatter) Format(timeString, levelString, prefix, caller string, line uint, msg string, data Data) []byte {
	var buf bytes.Buffer
	if data != nil {
		msg = appendData(msg, data)
	}
	for _, part := range f.parts {
		switch part {
		case "%t":
			buf.WriteString(timeString)
		case "%l":
		  buf.WriteString(levelString)
		case "%p":
			buf.WriteString(prefix)
		case "%c":
			buf.WriteString(caller)
		case "%n":
			buf.WriteString(strconv.FormatUint(uint64(line), 10))
		case "%m":
		  buf.WriteString(msg)
		case "%%":
			buf.WriteRune('%')
		default:
			buf.WriteString(part)
		}
	}

	return buf.Bytes()
}

func appendData(msg string, data Data) string {
	var buf bytes.Buffer
	buf.WriteString(msg)
	for key, value := range data {
		// TODO: faster, safer way of doing this
		buf.WriteString(fmt.Sprintf(" %s=%v", key, value))
	}
	return buf.String()
}
