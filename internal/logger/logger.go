package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(message string, v ...interface{}) {
	Logger.Printf("[INFO] "+message, v...)
}

func LogWarning(message string, v ...interface{}) {
	Logger.Printf("[WARNING] "+message, v...)
}

func LogError(message string, v ...interface{}) {
	Logger.Printf("[ERROR] "+message, v...)
}

func LogFatal(message string, v ...interface{}) {
	Logger.Fatalf("[FATAL] "+message, v...)
}

func LogPanic(message string, v ...interface{}) {
	Logger.Panicf("[PANIC] "+message, v...)
}
