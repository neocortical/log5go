package log5go

import (
	"os"
	"strings"
)

const (
	L5G_LOG_FILE_NAME   = "L5G_LOG_FILE_NAME"
	L5G_LOG_LEVEL       = "L5G_LOG_LEVEL"
	L5G_LOG_LINE_LENGTH = "L5G_LOG_LINE_LENGTH"
)

func GetLogger(key string) (l Log5Go) {
	l = GetOrCreate(key, func() (_ Log5Go) {
		return createLogFromEnvVars()
	}).WithPrefix(key)
	return
}

func createLogFromEnvVars() (l Log5Go) {
	l = getBaseLogWithLevel(os.Getenv(L5G_LOG_LEVEL))

	path, file := parseFilenameAndPath(os.Getenv(L5G_LOG_FILE_NAME))
	if path != "" && file != "" {
		l.ToFile(path, file)
	} else {
		// default to Stdout
		l.ToStdout().WithStderr()
	}

	switch os.Getenv(L5G_LOG_LINE_LENGTH) {
	case "NONE":
		// no-op
	case "LONG":
		l.WithLongLines()
	case "SHORT":
		l.WithShortLines()
	case "":
		l.WithShortLines()
	}

	return
}

func getBaseLogWithLevel(loglevel string) (l Log5Go) {
	for k, v := range levelMap {
		if v == loglevel {
			l = Logger(k)
			return
		}
	}

	l = Logger(LogAll)
	return
}

func parseFilenameAndPath(fullname string) (path, filename string) {
	if "" == fullname {
		return
	}

	i := strings.LastIndex(fullname, "/")

	if i+1 == len(fullname) {
		return
	}

	path = fullname[:i]
	filename = fullname[i+1:]
	return
}
