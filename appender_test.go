package log5go

import (
	"testing"
)

var appenderTests = map[string]string{
	"":      "\n",
	"\n":    "\n",
	"hello": "hello\n",
}

func TestTerminateWithNewline(t *testing.T) {
	for input, expected := range appenderTests {
		var buf []byte
		buf = append(buf, input...)
		TerminateMessageWithNewline(&buf)
		if string(buf) != expected {
			t.Errorf("expected %s but got %s", expected, string(buf))
		}
	}
}
