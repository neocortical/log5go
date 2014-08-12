package log5go

import (
	"testing"
)

func TestPutAndGet(t *testing.T) {
	log, _ := Log(LogAll).ToStdout().Build()

	err := loggerRegistry.Put("foo", log)
	if err != nil {
		t.Errorf("error putting logger into registry: %v", err)
	}

	log, err = loggerRegistry.Get("foo")
	if err != nil {
		t.Errorf("error getting logger from registry: %v", err)
	}

	err = loggerRegistry.Put("foo", log)
	if err == nil {
		t.Error("expected error when putting logger in registry with duplicate key")
	}
}
