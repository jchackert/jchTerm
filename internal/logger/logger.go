package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	logFile *os.File
	logger  *log.Logger
)

// Init initializes the logger with the given log file path
func Init(logFilePath string) error {
	var err error

	// Ensure the directory exists
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open the log file
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Create a new logger
	logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return nil
}

// Log writes a log message with the given format and arguments
func Log(format string, v ...interface{}) {
	if logger == nil {
		return
	}

	// Get the caller's file and line number
	_, file, line, _ := runtime.Caller(1)

	// Format the message
	message := fmt.Sprintf(format, v...)
	logMessage := fmt.Sprintf("%s:%d: %s", filepath.Base(file), line, message)

	// Write the log message
	logger.Println(logMessage)
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// GetLogFilePath returns a log file path based on the current date
func GetLogFilePath(baseDir string) string {
	now := time.Now()
	fileName := fmt.Sprintf("jchterm_%s.log", now.Format("2006-01-02"))
	return filepath.Join(baseDir, fileName)
}
