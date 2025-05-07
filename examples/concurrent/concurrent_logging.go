package main

import (
	"fmt"
	"github.com/Fernando2706/howl"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// simulateWork simulates a task that takes some time to complete
func simulateWork(id int) (int, time.Duration) {
	// Simulate variable work duration
	duration := time.Duration(100+rand.Intn(900)) * time.Millisecond
	time.Sleep(duration)
	
	// Simulate a result
	result := rand.Intn(1000)
	
	return result, duration
}

// worker represents a worker that performs tasks and logs its activity
func worker(id int, logger *howl.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Create a worker-specific logger with the worker ID as a field
	workerLogger := logger.WithField("worker_id", id)
	
	// Log worker start
	workerLogger.Info("Worker started")
	
	// Perform 5 tasks
	for i := 1; i <= 5; i++ {
		// Create a task-specific logger
		taskLogger := workerLogger.WithFields(map[string]interface{}{
			"task_id":    fmt.Sprintf("%d-%d", id, i),
			"task_count": i,
		})
		
		taskLogger.Debug("Starting task")
		
		// Simulate the task
		result, duration := simulateWork(id)
		
		// Log the task result with different levels based on the result
		if result < 200 {
			taskLogger.WithFields(map[string]interface{}{
				"result":   result,
				"duration": duration.Milliseconds(),
			}).Error("Task failed with low value")
		} else if result < 700 {
			taskLogger.WithFields(map[string]interface{}{
				"result":   result,
				"duration": duration.Milliseconds(),
			}).Warn("Task completed with average value")
		} else {
			taskLogger.WithFields(map[string]interface{}{
				"result":   result,
				"duration": duration.Milliseconds(),
				"metrics": map[string]interface{}{
					"efficiency": float64(result) / float64(duration.Milliseconds()),
					"quality":    float64(result) / 1000.0,
				},
			}).Info("Task completed successfully with high value")
		}
		
		// Simulate random errors
		if rand.Intn(10) < 2 { // 20% chance of error
			err := fmt.Errorf("random error in task %d-%d", id, i)
			taskLogger.WithError(err).Error("Encountered an error during task execution")
		}
	}
	
	// Log worker completion
	workerLogger.WithField("tasks_completed", 5).Info("Worker finished all tasks")
}

func main() {
	// Create logs directory if it doesn't exist
	logsDir := "examples/logs"
	err := os.MkdirAll(logsDir, 0755)
	if err != nil {
		fmt.Printf("Error creating logs directory: %v\n", err)
	}
	
	// Create a file for logging
	logFile, err := os.Create(filepath.Join(logsDir, "concurrent.log"))
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
	} else {
		defer logFile.Close()
	}
	
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Create a logger with JSON formatting
	logger := howl.New(
		howl.DefaultConfig().
			WithJSON(true).
			WithTimestamp(true).
			WithLevel(howl.DebugLevel),
	)
	
	// Create a file logger for saving logs to file
	var fileLogger *howl.Logger
	if logFile != nil {
		fileLogger = howl.New(
			howl.DefaultConfig().
				WithJSON(true).
				WithColor(false).
				WithTimestamp(true).
				WithLevel(howl.DebugLevel),
		)
		fileLogger.SetOutput(logFile)
		fileLogger.Info("Concurrent logging started")
	}
	
	// Log application start
	logger.WithFields(map[string]interface{}{
		"app_name":    "Concurrent Logger Example",
		"version":     "1.0.0",
		"start_time":  time.Now().Format(time.RFC3339),
		"num_workers": 5,
	}).Info("Application started")
	
	// Create a wait group to wait for all workers to complete
	var wg sync.WaitGroup
	
	// Start multiple workers concurrently
	numWorkers := 5
	wg.Add(numWorkers)
	
	logger.Info("Starting workers")
	
	for i := 1; i <= numWorkers; i++ {
		go worker(i, logger, &wg)
	}
	
	// Wait for all workers to complete
	wg.Wait()
	
	// Log application completion
	logger.WithField("end_time", time.Now().Format(time.RFC3339)).Info("All workers completed, application shutting down")
	
	// Log to file if available
	if fileLogger != nil {
		fileLogger.WithField("end_time", time.Now().Format(time.RFC3339)).Info("All workers completed, application shutting down")
	}
	
	// Demonstrate concurrent access to the same logger from multiple goroutines
	logger.Info("Demonstrating concurrent logging from multiple goroutines")
	
	// Create a wait group for the concurrent logging demonstration
	var wg2 sync.WaitGroup
	wg2.Add(100)
	
	// Log concurrently from 100 goroutines
	for i := 1; i <= 100; i++ {
		go func(id int) {
			defer wg2.Done()
			logger.WithFields(map[string]interface{}{
				"goroutine_id": id,
				"timestamp":    time.Now().UnixNano(),
			}).Info("Concurrent log entry")
		}(i)
	}
	
	// Wait for all concurrent logs to complete
	wg2.Wait()
	
	logger.Info("Concurrent logging demonstration completed")
	
	// Log to file if available
	if fileLogger != nil {
		fileLogger.Info("Concurrent logging demonstration completed")
		fmt.Printf("Logs have been saved to %s\n", filepath.Join(logsDir, "concurrent.log"))
	}
}
