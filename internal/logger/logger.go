package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

var logLevel = "INFO"

func InitLogger(level string) {
	logLevel = level
	Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func shouldLog(level string) bool {
	levels := map[string]int{"DEBUG": 1, "INFO": 2, "WARNING": 3, "ERROR": 4, "FATAL": 5}
	return levels[level] >= levels[logLevel]
}

func LogInfo(message string, v ...interface{}) {
	if shouldLog("INFO") {
		Logger.Printf(green+"[INFO] "+reset+message, v...)
	}
}

func LogWarning(message string, v ...interface{}) {
	if shouldLog("WARNING") {
		Logger.Printf(magenta+"[WARNING] "+reset+message, v...)
	}
}

func LogDebug(message string, v ...interface{}) {
	if shouldLog("DEBUG") {
		Logger.Printf(yellow+"[DEBUG] "+reset+message, v...)
	}
}

func LogError(message string, v ...interface{}) {
	if shouldLog("ERROR") {
		Logger.Printf(red+"[ERROR] "+message+reset, v...)
	}
}

func LogFatal(message string, v ...interface{}) {
	if shouldLog("FATAL") {
		Logger.Fatalf(blue+"[FATAL] "+message+reset, v...)
	}
}
