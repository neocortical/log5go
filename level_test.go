package log5go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldHaveTenInitialLevels(t *testing.T) {
	assert.Equal(t, 10, len(levelMap))
}

func Test_GetLogLevelString(t *testing.T) {
	name := GetLogLevelString(LogInfo)
	assert.Equal(t, "INFO", name)

	// bad log level
	name = GetLogLevelString(LogLevel(666))
	assert.Equal(t, "", name)
}

func Test_ShouldRegisterLogLevels(t *testing.T) {

	RegisterLogLevel(LogInfo, "HEY")
	if levelMap[LogInfo] != "HEY" {
		t.Errorf("expected HEY but got %s", levelMap[LogInfo])
	}

	name := GetLogLevelString(LogInfo)
	if name != "HEY" {
		t.Errorf("expected HEY but got %s", name)
	}

	// cleanup
	RegisterLogLevel(LogInfo, "INFO")

	customLevel := LogWarn + 1
	name = GetLogLevelString(customLevel)
	if name != "" {
		t.Errorf("expected empty string but got %s", name)
	}

	RegisterLogLevel(customLevel, "CUSTOM_WARN")
	if levelMap[customLevel] != "CUSTOM_WARN" {
		t.Errorf("expected CUSTOM_WARN but got %s", levelMap[customLevel])
	}

	name = GetLogLevelString(customLevel)
	if name != "CUSTOM_WARN" {
		t.Errorf("expected CUSTOM_WARN but got %s", name)
	}

	DeregisterLogLevel(customLevel)
	if levelMap[customLevel] != "" {
		t.Errorf("expected empty string but got %s", levelMap[customLevel])
	}

	name = GetLogLevelString(customLevel)
	if name != "" {
		t.Errorf("expected empty string but got %s", name)
	}
}
