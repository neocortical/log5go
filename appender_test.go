package log5go

import (
	"testing"
)

var appenderTests = map[string]string {
	"": "\n",
	"\n": "\n",
	"hello": "hello\n",
}

func TestTerminateWithNewline(t *testing.T) {
	for input, expected := range appenderTests {
		actual := TerminateMessageWithNewline([]byte(input))
		if string(actual) != expected {
			t.Errorf("expected %s but got %s", expected, actual)
		}
	}
}
