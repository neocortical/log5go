package log5go

import (
	"bytes"
	"testing"
)

var boundLoggerTests = []loggerTest {
	{"withdata", "hello, %s", []interface{}{"world"}, Data{"foo":"bar"}, TF_GoStd, "prefix", nil, "^" + Rxdate + " " + Rxtime + " {{level}} prefix: hello, world foo=\"bar\"\n$"},
}

func TestWithData(t *testing.T) {
	inner, _ := Logger(LogAll).(*logger)
	bl := &boundLogger{l: inner, data:make(Data)}

	if len(bl.data) != 0 {
		t.Errorf("expected empty data but got %v", bl.data)
	}

	bl, ok := bl.WithData(Data{"foo":"bar", "baz":1}).(*boundLogger)
	if !ok {
		t.Error("expected *boundLogger back from WithData")
	}
	if len(bl.data) != 2 || bl.data["foo"] != "bar" || bl.data["baz"] != 1 {
		t.Errorf("expected two-element data but got %v", bl.data)
	}

	// piling on
	bl, ok = bl.WithData(Data{"baz":2, "qux":3.14}).(*boundLogger)
	if !ok {
		t.Error("expected *boundLogger back from WithData")
	}
	if len(bl.data) != 3 || bl.data["foo"] != "bar" || bl.data["baz"] != 2 || bl.data["qux"] != 3.14 {
		t.Errorf("expected three-element data but got %v", bl.data)
	}
}

func TestAllBoundLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	appender := &writerAppender{dest:&buf}

	for _, test := range loggerTests {
		l := &logger{
			level: LogAll,
			formatter: test.formatter,
			appender: appender,
			timeFormat: test.timeFormat,
			prefix: test.prefix,
		}
		bl := &boundLogger{
			l: l,
			data: test.data,
		}
		runLoggerLevelTest(bl, &buf, &test, t)
	}

	for _, test := range boundLoggerTests {
		l := &logger{
			level: LogAll,
			formatter: test.formatter,
			appender: appender,
			timeFormat: test.timeFormat,
			prefix: test.prefix,
		}
		bl := &boundLogger{
			l: l,
			data: test.data,
		}
		runLoggerLevelTest(bl, &buf, &test, t)
	}
}
