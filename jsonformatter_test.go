package log5go

import (
	"testing"
	"time"
)

func TestDefaultJsonFormatter(t *testing.T) {
	d := Data{
		"bar": "baz",
	}
	var buf []byte

	theTime := time.Unix(1423343766, 0)
	jsonFormatter := &jsonFormatter{timeFormat: TF_GoStd, lines: true}

	jsonFormatter.Format(theTime, LogInfo, "prefix", "acme.go", 123, "foo", d, &buf)
	expected := "{\"time\":\"" + theTime.Format(TF_GoStd) + "\",\"level\":\"INFO\",\"prefix\":\"prefix\",\"line\":\"acme.go:123\",\"msg\":\"foo\",\"data\":{\"bar\":\"baz\"}}"
	if string(buf) != expected {
		t.Errorf("expected \n%s\n  but got \n%s", expected, string(buf))
	}
}
