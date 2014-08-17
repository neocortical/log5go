package log5go

import (
	"testing"
)

func TestRegisterLogLevel(t *testing.T) {
	if len(levelMap) != 7 {
		t.Errorf("error map not initialized to correct size")
	}

	name := GetLogLevelString(LogInfo)
	if name != "INFO" {
		t.Errorf("expected INFO but got %s", name)
	}

	RegisterLogLevel(LogInfo, "HEY")
	if levelMap[LogInfo] != "HEY" {
		t.Errorf("expected HEY but got %s", levelMap[LogInfo])
	}

	name = GetLogLevelString(LogInfo)
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
