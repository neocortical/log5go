package log5go

import (
	"os"
	"strconv"
	"strings"
)

const (
	L5G_LOG_FILE_NAME                  = "L5G_LOG_FILE_NAME"
	L5G_LOG_LEVEL                      = "L5G_LOG_LEVEL"
	L5G_LOG_LINE_LENGTH                = "L5G_LOG_LINE_LENGTH"
	L5G_LOG_FILE_ROTATION_FREQUENCY    = "L5G_LOG_FILE_ROTATION_FREQUENCY"
	L5G_LOG_FILE_ROTATION_KEEP_N_FILES = "L5G_LOG_FILE_ROTATION_KEEP_N_FILES"
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

		n := parseKeepNFilesInt(os.Getenv(L5G_LOG_FILE_ROTATION_KEEP_N_FILES))
		// file rotation
		switch os.Getenv(L5G_LOG_FILE_ROTATION_FREQUENCY) {
		case "MINUTE":
			l.WithRotation(RollMinutely, n)
		case "HOUR":
			l.WithRotation(RollHourly, n)
		case "DAY":
			l.WithRotation(RollDaily, n)
		case "WEEK":
			l.WithRotation(RollWeekly, n)
		}

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

func parseKeepNFilesInt(str string) (i int) {
	i, err := strconv.Atoi(str)
	if err != nil {
		// default to 1
		i = 1
	}
	return
}
