package utils

import (
	"log"
	"os"
	"runtime"
)

// GetStandardLogger returns a standard logger for server error logging
func GetStandardLogger() *log.Logger {
	return log.New(os.Stderr, "[SERVER] ", log.LstdFlags|log.Lshortfile)
}

func Logger(logLevel string, message string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}
	// FIXME limit the show logLevel if debug is not considered, then don't
	switch logLevel {
	case "debug":
		log.Printf("[DEBUG] %s:%d - %s\n", file, line, message)
		return
	case "info":
		log.Printf("[INFO] %s:%d - %s\n", file, line, message)
		return
	case "warn":
		log.Printf("[WARN] %s:%d - %s\n", file, line, message)
		return
	case "err":
		log.Printf("[ERROR] %s:%d - %s\n", file, line, message)
		return
	case "fatal":
		log.Fatalf("[FATAL] %s:%d - %s\n", file, line, message)
	}
}
