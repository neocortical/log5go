package log5go

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	L5G_LOG_FILE_NAME                  = "L5G_LOG_FILE_NAME"
	L5G_LOG_LEVEL                      = "L5G_LOG_LEVEL"
	L5G_LOG_LINE_LENGTH                = "L5G_LOG_LINE_LENGTH"
	L5G_LOG_FILE_ROTATION_FREQUENCY    = "L5G_LOG_FILE_ROTATION_FREQUENCY"
	L5G_LOG_FILE_ROTATION_KEEP_N_FILES = "L5G_LOG_FILE_ROTATION_KEEP_N_FILES"
)

type logconf struct {
	logLevel             LogLevel
	logFilePath          string
	logFileName          string
	logFileRollFrequency rollFrequency
	keepNFiles           int
	logLineLength        string
}

var lock = sync.Mutex{}
var conf *logconf

func GetConsoleLogger() (l Log5Go) {
	l = GetOrCreate("console", func() (_ Log5Go) {
		return Logger(LogAll).ToStdout()
	})
	return
}

func GetLogger(key string) (l Log5Go) {
	l = GetOrCreate(key, func() (_ Log5Go) {
		return createLogFromEnvVars()
	}).WithPrefix(key)
	return
}

func loadEnv() {
	lock.Lock()
	defer lock.Unlock()

	if conf == nil {
		console := GetConsoleLogger()
		console.Info("Initializing Log5Go ... ")

		conf = &logconf{
			logLevel: parseLogLevel(os.Getenv(L5G_LOG_LEVEL)),
		}

		console.Info("[Log Level: %v]", levelMap[conf.logLevel])

		conf.logFilePath, conf.logFileName = parseFilenameAndPath(os.Getenv(L5G_LOG_FILE_NAME))
		if conf.logFilePath != "" && conf.logFileName != "" {
			console.Info("[Logging to file: %s/%s]", conf.logFilePath, conf.logFileName)
			conf.keepNFiles = parseKeepNFilesInt(os.Getenv(L5G_LOG_FILE_ROTATION_KEEP_N_FILES))
			logFileRollFrequency, rollLabel := parseFileRotationFrequency(os.Getenv(L5G_LOG_FILE_ROTATION_FREQUENCY))
			conf.logFileRollFrequency = logFileRollFrequency

			console.Info("[Logfile Roll Frequency: %s]", rollLabel)
			console.Info("[Logfiles to keep: %d]", conf.keepNFiles)

		} else {
			console.Info("[Logging to Stdout]")
		}

		conf.logLineLength = parseLogLineLength(os.Getenv(L5G_LOG_LINE_LENGTH))
		console.Info("[Log Line Length: %s]", conf.logLineLength)
	}
}

func createLogFromEnvVars() (l Log5Go) {
	loadEnv()

	// create a base Log5Go with the appropriate log level
	l = Logger(conf.logLevel)

	if conf.logFilePath != "" && conf.logFileName != "" {
		l.ToFile(conf.logFilePath, conf.logFileName)

		// file rotation
		if conf.logFileRollFrequency != RollNone {
			l.WithRotation(conf.logFileRollFrequency, conf.keepNFiles)
		}
	} else {
		// default to Stdout
		l.ToStdout().WithStderr()
	}

	switch conf.logLineLength {
	case "LONG":
		l.WithLongLines()
	case "SHORT", "NONE":
		l.WithShortLines()
	}

	return
}

func parseLogLineLength(token string) (str string) {
	str = "NONE"
	switch token {
	case "LONG", "SHORT":
		str = token
	}
	return str
}

func parseLogLevel(logLevelStr string) (level LogLevel) {
	for k, v := range levelMap {
		if v == logLevelStr {
			level = k
			return
		}
	}

	level = LogAll
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

func parseFileRotationFrequency(str string) (freq rollFrequency, label string) {
	label = str
	switch str {
	case "MINUTE":
		freq = RollMinutely
		return
	case "HOUR":
		freq = RollHourly
		return
	case "DAY":
		freq = RollDaily
		return
	case "WEEK":
		freq = RollWeekly
		return
	}

	return RollNone, "NONE"
}

func parseKeepNFilesInt(str string) (i int) {
	i, err := strconv.Atoi(str)
	if err != nil {
		// default to 1
		i = 1
	}
	return
}
