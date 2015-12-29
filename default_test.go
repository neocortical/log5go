package log5go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseLogLevel_Default(t *testing.T) {
	assert.Equal(t, LogAll, parseLogLevel(""))
}

func Test_parseLogLevel_LogNotice(t *testing.T) {
	assert.Equal(t, LogNotice, parseLogLevel("NOTICE"))
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

func Test_parseLogLineLength_Default(t *testing.T) {
	assert.Equal(t, "NONE", parseLogLineLength(""))
}

func Test_parseLogLineLength_Bad(t *testing.T) {
	assert.Equal(t, "NONE", parseLogLineLength("FOO"))
}

func Test_parseLogLineLength_None(t *testing.T) {
	assert.Equal(t, "NONE", parseLogLineLength("NONE"))
}

func Test_parseLogLineLength_Long(t *testing.T) {
	assert.Equal(t, "LONG", parseLogLineLength("LONG"))
}

func Test_parseLogLineLength_Short(t *testing.T) {
	assert.Equal(t, "SHORT", parseLogLineLength("SHORT"))
}
