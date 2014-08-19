package log5go

import (
	"testing"
)

func TestDefaultJsonFormatter(t *testing.T) {
	d := Data{
		"bar": "baz",
	}
	var buf []byte

	defaultJsonFormatter.Format("time", "level", "prefix", "acme.go", 123, "foo", d, &buf)
	expected := "{\"time\":\"time\",\"level\":\"level\",\"prefix\":\"prefix\",\"line\":\"acme.go:123\",\"msg\":\"foo\",\"data\":{\"bar\":\"baz\"}}"
	if string(buf) != expected {
		t.Errorf("expected \n%s\n  but got \n%s", expected, string(buf))
	}
}
