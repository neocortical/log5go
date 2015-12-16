package log5go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getBaseLogWithLevel_Default(t *testing.T) {
	l := getBaseLogWithLevel("")
	if assert.NotNil(t, l) {
		assert.Equal(t, LogAll, l.LogLevel())
	}
}

func Test_getBaseLogWithLevel_LogNotice(t *testing.T) {
	l := getBaseLogWithLevel("NOTICE")
	if assert.NotNil(t, l) {
		assert.Equal(t, LogNotice, l.LogLevel())
	}
}

func Test_parseFilenameAndPath_Empty(t *testing.T) {
	path, file := parseFilenameAndPath("")
	assert.Equal(t, "", path)
	assert.Equal(t, "", file)
}

func Test_parseFilenameAndPath_HappyPath(t *testing.T) {
	path, file := parseFilenameAndPath("/foo/1/bar")
	assert.Equal(t, "/foo/1", path)
	assert.Equal(t, "bar", file)
}

func Test_parseFilenameAndPath_TraliningSlash(t *testing.T) {
	path, file := parseFilenameAndPath("/foo/bar/")
	assert.Equal(t, "", path)
	assert.Equal(t, "", file)
}

func Test_parseKeepNFilesInt_Empty(t *testing.T) {
	assert.Equal(t, 1, parseKeepNFilesInt(""))
}

func Test_parseKeepNFilesInt_10(t *testing.T) {
	assert.Equal(t, 10, parseKeepNFilesInt("10"))
}

func Test_parseKeepNFilesInt_1(t *testing.T) {
	assert.Equal(t, 1, parseKeepNFilesInt("1"))
}

func Test_parseKeepNFilesInt_NotNumeric(t *testing.T) {
	assert.Equal(t, 1, parseKeepNFilesInt("foobar"))
}
