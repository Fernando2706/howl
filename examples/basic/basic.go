package main

import (
	"errors"
	"os"
	"time"

	"github.com/Fernando2706/howl"
)

func main() {
	// PART 1: Basic logging examples
	basicLoggingExamples()

	// PART 2: JSON logging examples
	jsonLoggingExamples()

	// PART 3: File logging examples
	fileLoggingExamples()
}

func basicLoggingExamples() {
	// Create a logger with default configuration
	logger := howl.Default()
	logger.Info("Application started")

	// Log messages with different levels
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Log with structured fields
	logger.WithField("user", "admin").Info("User logged in")

	logger.WithFields(map[string]interface{}{
		"user": "admin",
		"ip":   "192.168.1.1",
	}).Info("Connection details")

	// Log with error
	err := errors.New("something went wrong")
	logger.WithError(err).Error("Operation failed")
}

func jsonLoggingExamples() {
	// Create a logger with JSON formatting enabled
	jsonLogger := howl.New(
		howl.DefaultConfig().
			WithJSON(true).
			WithColor(true).
			WithTimestamp(true),
	)

	// Basic JSON log
	jsonLogger.Info("JSON logging started")

	// JSON log with a single field
	jsonLogger.WithField("version", "1.0.0").Info("Version information")

	// JSON log with multiple fields
	jsonLogger.WithFields(map[string]interface{}{
		"user":      "admin",
		"ip":        "192.168.1.1",
		"timestamp": time.Now().Unix(),
	}).Info("User login")

	// JSON log with nested data
	userData := map[string]interface{}{
		"id":       12345,
		"username": "admin",
		"roles":    []string{"admin", "user"},
		"settings": map[string]interface{}{
			"theme":       "dark",
			"language":    "en",
			"preferences": map[string]bool{"notifications": true, "sounds": false},
		},
	}
	jsonLogger.WithField("user_data", userData).Info("User profile loaded")

	// JSON log with request information
	jsonLogger.WithFields(map[string]interface{}{
		"request_id": "abc-123-xyz",
		"method":     "GET",
		"path":       "/api/users",
		"status":     200,
		"duration":   157.5,
	}).Info("API request completed")

	// JSON log with error information
	jsonLogger.WithFields(map[string]interface{}{
		"error_code": 500,
		"component":  "database",
		"query":      "SELECT * FROM users",
		"stack": []string{
			"main.connectDB()",
			"main.initializeApp()",
			"main.main()",
		},
	}).Error("Database connection failed")

	// JSON log with application metrics
	jsonLogger.WithFields(map[string]interface{}{
		"memory_usage": 1024.5,
		"cpu_usage":    45.2,
		"goroutines":   15,
		"uptime":       3600,
		"metrics": map[string]float64{
			"requests_per_second": 152.3,
			"average_response_ms": 45.7,
			"error_rate":          0.01,
		},
	}).Info("Application metrics")
}

func fileLoggingExamples() {
	// Create a regular logger for console output
	logger := howl.Default()

	// Create a file for logging
	file, err := os.Create("examples/logs/app.log")
	if err != nil {
		logger.WithError(err).Fatal("Could not create log file")
	}
	defer file.Close()

	// Create a logger that writes to the file
	fileLogger := howl.New(howl.DefaultConfig())
	fileLogger.SetOutput(file)

	// Log to the file
	fileLogger.Info("This message goes to the file")
	fileLogger.WithField("source", "file_logger").Info("Log with field")

	// Create a JSON logger that writes to a file
	jsonFileLogger := howl.New(
		howl.DefaultConfig().
			WithJSON(true).
			WithColor(false).
			WithTimestamp(true),
	)

	// Create a file for JSON logging
	jsonFile, err := os.Create("examples/logs/json_logs.log")
	if err != nil {
		logger.WithError(err).Fatal("Could not create JSON log file")
	}
	defer jsonFile.Close()

	jsonFileLogger.SetOutput(jsonFile)

	// Log JSON to the file
	jsonFileLogger.Info("JSON logging to file started")
	jsonFileLogger.WithFields(map[string]interface{}{
		"app":     "example",
		"version": "1.0.0",
		"env":     "development",
	}).Info("Application configuration")

	// Show in console that we've written to the files
	logger.WithField("file", "examples/logs/app.log").Info("Logs written to file")
	logger.WithField("file", "examples/logs/json_logs.log").Info("JSON logs written to file")
}
