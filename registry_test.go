package log5go

import (
	"testing"
)

func TestPutAndGet(t *testing.T) {
	log := Logger(LogAll).ToStdout()

	loggerRegistry.Put("foo", log)

	log, err := loggerRegistry.Get("foo")
	if err != nil {
		t.Errorf("error getting logger from registry: %v", err)
	}
}
