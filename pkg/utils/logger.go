package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity of the log message
type LogLevel string

const (
	// LogLevelDebug represents debug level logs
	LogLevelDebug LogLevel = "DEBUG"
	// LogLevelInfo represents info level logs
	LogLevelInfo LogLevel = "INFO"
	// LogLevelWarn represents warning level logs
	LogLevelWarn LogLevel = "WARN"
	// LogLevelError represents error level logs
	LogLevelError LogLevel = "ERROR"
	// LogLevelFatal represents fatal level logs
	LogLevelFatal LogLevel = "FATAL"
)

var (
	standardLogger *log.Logger
)

func init() {
	standardLogger = log.New(os.Stdout, "", log.LstdFlags)
}

// Logger logs a message with the specified level and message
func Logger(level string, message string) {
	logLevel := LogLevel(level)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Format: [TIMESTAMP] [LEVEL] Message
	logMessage := fmt.Sprintf("[%s] [%s] %s", timestamp, logLevel, message)

	switch logLevel {
	case LogLevelFatal:
		standardLogger.Fatal(logMessage)
	case LogLevelError:
		standardLogger.Print(logMessage)
	case LogLevelWarn:
		standardLogger.Print(logMessage)
	case LogLevelInfo:
		standardLogger.Print(logMessage)
	case LogLevelDebug:
		standardLogger.Print(logMessage)
	default:
		standardLogger.Printf("[%s] [UNKNOWN] %s", timestamp, message)
	}
}

// Debug logs a debug level message
func Debug(message string) {
	Logger(string(LogLevelDebug), message)
}

// Info logs an info level message
func Info(message string) {
	Logger(string(LogLevelInfo), message)
}

// Warn logs a warning level message
func Warn(message string) {
	Logger(string(LogLevelWarn), message)
}

// Error logs an error level message
func Error(message string) {
	Logger(string(LogLevelError), message)
}

// Fatal logs a fatal level message and exits the program
func Fatal(message string) {
	Logger(string(LogLevelFatal), message)
}

// Debugf logs a formatted debug level message
func Debugf(format string, args ...interface{}) {
	Debug(fmt.Sprintf(format, args...))
}

// Infof logs a formatted info level message
func Infof(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning level message
func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error level message
func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...))
}

// Fatalf logs a formatted fatal level message and exits the program
func Fatalf(format string, args ...interface{}) {
	Fatal(fmt.Sprintf(format, args...))
}

// GetStandardLogger returns the standard logger instance
func GetStandardLogger() *log.Logger {
	return standardLogger
}
