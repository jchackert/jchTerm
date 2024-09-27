package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jchackert/jchterm/internal/logger"
	"github.com/jchackert/jchterm/internal/tui"
)

func main() {
	// Get the directory of the executable
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error getting executable directory: %v\n", err)
		os.Exit(1)
	}

	// Get the project root directory (assuming the executable is in the project root)
	projectRoot := execDir

	// Construct the path to the logs directory
	logsDir := filepath.Join(projectRoot, "logs")

	// Ensure the logs directory exists
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		fmt.Printf("Error creating logs directory: %v\n", err)
		os.Exit(1)
	}

	// Get the log file path
	logFile := logger.GetLogFilePath(logsDir)

	// Initialize the logger
	if err := logger.Init(logFile); err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Log application start and configuration
	logger.Log("Application started")
	logger.Log("Executable directory: %s", execDir)
	logger.Log("Project root directory: %s", projectRoot)
	logger.Log("Logs directory: %s", logsDir)
	logger.Log("Log file path: %s", logFile)

	// Run the TUI
	if err := tui.Run(); err != nil {
		logger.Log("Error running TUI: %v", err)
		fmt.Printf("An error occurred. Please check the log file at %s for details.\n", logFile)
		os.Exit(1)
	}

	// Log application exit
	logger.Log("Application exiting normally")
}
