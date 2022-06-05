package archivist

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var logFlags = [5]bool{false, true, true, true, true}
var logLevels = map[string]int{"debug": 0, "info": 1, "warning": 2, "error": 3, "fatal": 4}

var logger = log.New(os.Stdout, "", 0)

func Init(logLevel string, logTarget string, arg ...string) {
	// first we gonne dispatch the LOG_LEVEL default is info
	SetLogLevel(logLevel)

	// first we gonne check the log target
	// and may change the target based on config
	switch logTarget {
	case "file":
		Info("Setting logger output to file")
		// make sure we have a given logPath that isnt empty
		// we assume the 3rd param (1st arg) is logPath based
		// on logTarget file
		if 0 < len(arg) && "" != arg[0] {
			logPath := arg[0]
			if _, err := os.Stat(logPath); errors.Is(err, os.ErrNotExist) {
				file, err := os.Create(logPath)
				if err != nil {
					Error("Could not create new logFile on given LOG_PATH ", logPath)
					break
				} else {
					Info("Created new logFile on given LOG_PATH ", logPath)
				}
				logger.SetOutput(file)
			} else {
				file, err := os.Open(logPath) // For read access.
				if nil != err {
					Error("Cannot open specified LOG_FILE", logPath)
				} else {
					logger.SetOutput(file)
				}
			}
		} else {
			Error("Invalid logPath supplied. Fallback to default", arg)
		}
	case "stdout":
		// since stdout is default we got nuttin to do here
		// keepin it for the case the default changes and we
		// have to init a stdout logger
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		// on default we gonne fallback to stdout aka default
		// since stdout is default we got nuttin to do here
		// keepin it for the case the default changes and we
		// have to init a stdout logger
	}
}

func store(message string, stype string, dump bool, formatted bool, params ...interface{}) {
	// dispatch the caller file+line number
	_, file, line, _ := runtime.Caller(2)
	// - - - - - - - - - - - - - - - - - - - - - - -
	// KEEP THIS #### THINKING ABOUT ADDING METHOD CALLER NAME
	// handle := runtime.FuncForPC(fn) // where fn is the first param given from runtime.Caller
	// fName := handle.Name()
	// - - - - - - - - - - - - - - - - - - - - - - -
	arrPackagePath := strings.Split(file, "/")
	packageFile := arrPackagePath[len(arrPackagePath)-1]

	// build the actual logline
	logLine := time.Now().Format("2006-01-02 15:04:05") + "|" + stype + "|" + packageFile + "#" + strconv.Itoa(line) + "|"
	if true == dump {
		if true == formatted {
			logLine = logLine + fmt.Sprintf(message, params)
		} else {
			logLine = logLine + message + "|" + fmt.Sprintf("%+v", params)
		}
	} else {
		logLine = logLine + message
	}

	// finally print it via logger (threadsafe)
	logger.Print(logLine)
}

func Error(message string, params ...interface{}) {
	// 3 = "error"
	if logFlags[3] {
		if 0 == len(params) {
			store(message, "error", false, false, "")
		} else {
			store(message, "error", true, false, params)
		}
	}
}

func ErrorF(message string, params ...interface{}) {
	// 3 = "error"
	if logFlags[3] {
		store(message, "error", true, true, params)
	}
}

func Fatal(message string, params ...interface{}) {
	// 4 = "fatal"
	if logFlags[4] {
		if 0 == len(params) {
			store(message, "fatal", false, false, "")
		} else {
			store(message, "fatal", true, false, params)
		}
	}
}

func FatalF(message string, params ...interface{}) {
	// 4 = "fatal"
	if logFlags[4] {
		store(message, "fatal", true, true, params)
	}
}

func Info(message string, params ...interface{}) {
	// 1 == "info"
	if logFlags[1] {
		if 0 == len(params) {
			store(message, "info", false, false, "")
		} else {
			store(message, "info", true, false, params)
		}
	}
}

func InfoF(message string, params ...interface{}) {
	// 1 = "info"
	if logFlags[1] {
		store(message, "info", true, true, params)
	}
}

func Warning(message string, params ...interface{}) {
	// 2 = "warning"
	if logFlags[2] {
		if 0 == len(params) {
			store(message, "warning", false, false, "")
		} else {
			store(message, "warning", true, false, params)
		}
	}
}

func WarningF(message string, params ...interface{}) {
	// 2 = "warning"
	if logFlags[2] {
		store(message, "warning", true, true, params)
	}
}

func Debug(message string, params ...interface{}) {
	// 0 = "debug"
	if logFlags[0] {
		if 0 == len(params) {
			store(message, "debug", false, false, "")
		} else {
			store(message, "debug", true, false, params)
		}
	}
}

func DebugF(message string, params ...interface{}) {
	// 0 = "debug"
	if logFlags[0] {
		store(message, "debug", true, true, params)
	}
}

func SetLogLevel(logLevel string) {
	if logLevelID, ok := logLevels[logLevel]; ok {
		Info("Setting log level to", logLevel, logLevelID)
		for index, _ := range logFlags {
			if logLevelID <= index {
				logFlags[index] = true
			} else {
				logFlags[index] = false
			}
		}
	} else {
		Error("Given LOG_LEVEL is unknown", logLevel)
	}
}
