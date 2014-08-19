package log5go

import (
	"encoding/json"
	"fmt"
)

type jsonFormatter struct{}

var defaultJsonFormatter = &jsonFormatter{}

type jsonLog struct {
	Time   string                 `json:"time"`
	Level  string                 `json:"level"`
	Prefix string                 `json:"prefix,omitempty"`
	Line   string                 `json:"line,omitempty"`
	Msg    string                 `json:"msg"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

func (f *jsonFormatter) Format(timeString, levelString, prefix, caller string, line uint, msg string, data Data, out *[]byte) {
	output := jsonLog{
		Time:   timeString,
		Level:  levelString,
		Prefix: prefix,
		Line:   formatLine(caller, line),
		Msg:    msg,
		Data:   data,
	}

	serialized, err := json.Marshal(output) // TODO: remove intermediate string
	if err == nil {
		*out = append(*out, serialized...)
	}
}

func formatLine(caller string, line uint) string {
	if caller == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", caller, line)
}
