package log5go

import (
	"testing"
	"time"
)

func TestNewStringFormatter(t *testing.T) {
	// all the stuff
	sf := NewStringFormatter("%t %l/%L %p (%c:%n): %艾未未 %m %d %% junk %")

	expected := []string{"%t", " ", "%l", "/L ", "%p", " (", "%c", ":", "%n", "): 艾未未 ", "%m", " d ", "%%", " junk "}

	if !testEq(sf.parts, expected) {
		t.Errorf("expected \n%v, but got \n%v", expected, sf.parts)
	}

	// empty format string
	sf = NewStringFormatter("")

	expected = []string{}

	if !testEq(sf.parts, expected) {
		t.Errorf("expected \n%v, but got \n%v", expected, sf.parts)
	}
}

func TestStringFormatterParse(t *testing.T) {
	theTime := time.Unix(1423343766, 0)

	var buf []byte
	sf := NewStringFormatter("%t %l %p (%c:%n): %m %%艾未未")
	sf.Format(theTime, LogInfo, "艾未未", "acme.go", 123, "hello?", nil, &buf)
	expected := theTime.Format(TF_GoStd) + " INFO 艾未未 (acme.go:123): hello? %艾未未"
	if expected != string(buf) {
		t.Errorf("expected %s but got %s", expected, string(buf))
	}

	buf = buf[:0]
	sf = NewStringFormatter("")
	sf.Format(theTime, LogInfo, "艾未未", "acme.go", 123, "hello?", nil, &buf)
	expected = ""
	if expected != string(buf) {
		t.Errorf("expected %s but got %s", expected, string(buf))
	}
}

func TestDataAppend(t *testing.T) {
	theTime := time.Unix(1423343766, 0)
	d := Data{
		"foo": "bar",
		"baz": 42,
	}

	var buf []byte

	sf := NewStringFormatter("%t %l %p: %m")
	sf.Format(theTime, LogInfo, "艾未未", "acme.go", 123, "hello?", d, &buf)
	expected := theTime.Format(TF_GoStd) + " INFO 艾未未: hello? foo=\"bar\" baz=42"
	expected2 := theTime.Format(TF_GoStd) + " INFO 艾未未: hello? baz=42 foo=\"bar\""
	if expected != string(buf) && expected2 != string(buf) {
		t.Errorf("expected %s but got %s", expected, string(buf))
	}
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
