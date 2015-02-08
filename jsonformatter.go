package log5go

import (
	"encoding/json"
	"fmt"
	"time"
)

type jsonFormatter struct {
	timeFormat string
	lines      bool
}

type jsonLog struct {
	Time   string                 `json:"time"`
	Level  string                 `json:"level"`
	Prefix string                 `json:"prefix,omitempty"`
	Line   string                 `json:"line,omitempty"`
	Msg    string                 `json:"msg"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

func (f *jsonFormatter) Format(tstamp time.Time, level LogLevel, prefix, caller string, line uint, msg string, data Data, out *[]byte) {
	output := jsonLog{
		Time:   tstamp.Format(f.timeFormat),
		Level:  GetLogLevelString(level),
		Prefix: prefix,
		Line:   f.formatLine(caller, line),
		Msg:    msg,
		Data:   data,
	}

	serialized, err := json.Marshal(output)
	if err == nil {
		*out = append(*out, serialized...)
	}
}

func (f *jsonFormatter) SetTimeFormat(timeFormat string) {
	f.timeFormat = timeFormat
}

func (f *jsonFormatter) SetLines(lines bool) {
	f.lines = lines
}

func (f *jsonFormatter) formatLine(caller string, line uint) string {
	if !f.lines || caller == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", caller, line)
}
