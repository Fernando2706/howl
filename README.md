# Howl

Howl is a flexible and powerful logging library for Go applications with a focus on simplicity and readability.

## Features

- Multiple log levels (Debug, Info, Warn, Error, Fatal)
- Support for different formats (text, JSON)
- Colorized output for improved readability
- Colorized JSON for better visualization of structured data
- Multiple output destinations (console, file)
- Structured logging with fields
- Context integration with extractors
- Thread-safe logging for concurrent environments
- Fluent and easy-to-use API

## Installation

```bash
go get github.com/Fernando2706/howl
```

## Basic Usage

```go
package main

import (
    "github.com/Fernando2706/howl"
)

func main() {
    // Create a logger with default configuration
    logger := howl.New(howl.DefaultConfig())
    
    // Simple logs
    logger.Debug("This is a debug message")
    logger.Info("Application started")
    logger.Warn("This is a warning")
    logger.Error("An error occurred")
    
    // Logs with structured fields
    logger.WithField("user", "admin").Info("User logged in")
    
    logger.WithFields(map[string]interface{}{
        "user": "admin",
        "ip":   "192.168.1.1",
    }).Info("Connection details")
    
    // Logs with errors
    err := errors.New("something went wrong")
    logger.WithError(err).Error("Operation failed")
}
```

## Custom Configuration

```go
package main

import (
    "os"
    "github.com/Fernando2706/howl"
)

func main() {
    // Create a logger with custom configuration
    logger := howl.New(
        howl.DefaultConfig().
            WithLevel(howl.DebugLevel).
            WithJSON(true).
            WithColor(true).
            WithTimestamp(true),
    )
    
    // Log to a file
    file, err := os.Create("app.log")
    if err != nil {
        logger.WithError(err).Fatal("Could not create log file")
    }
    defer file.Close()
    
    // Create a logger that writes to the file
    fileLogger := howl.New(howl.DefaultConfig())
    fileLogger.SetOutput(file)
    
    fileLogger.Info("This message goes to the file")
    logger.Info("This message goes to the console")
}
```

## Context Integration

```go
package main

import (
    "context"
    "github.com/Fernando2706/howl"
)

// RequestIDExtractor extracts the request ID from the context
func RequestIDExtractor(ctx context.Context) map[string]interface{} {
    fields := make(map[string]interface{})
    if requestID, ok := ctx.Value("request_id").(string); ok {
        fields["request_id"] = requestID
    }
    return fields
}

func main() {
    // Create a logger with context extractors
    logger := howl.New(
        howl.DefaultConfig().WithJSON(true),
    )
    
    // Add a context extractor
    logger = logger.WithContextExtractor(RequestIDExtractor)
    
    // Create a context with request ID
    ctx := context.Background()
    ctx = context.WithValue(ctx, "request_id", "req-123")
    
    // Log with context - automatically includes request_id
    logger.InfoContext(ctx, "Request started")
    
    // Log with additional fields
    logger.WithField("path", "/api/users").InfoContext(ctx, "Processing request")
}
```

## JSON Formatting

Howl supports colorized JSON output for better readability:

```go
logger := howl.New(
    howl.DefaultConfig().
        WithJSON(true).
        WithColor(true),
)

logger.WithFields(map[string]interface{}{
    "user": "admin",
    "user_data": map[string]interface{}{
        "id": 12345,
        "roles": []string{"admin", "user"},
        "settings": map[string]interface{}{
            "theme": "dark",
            "language": "en",
        },
    },
}).Info("User profile loaded")
```

## Thread Safety

Howl is designed to be thread-safe, making it suitable for concurrent applications:

```go
func worker(id int, logger *howl.Logger, wg *sync.WaitGroup) {
    defer wg.Done()
    
    logger.WithField("worker_id", id).Info("Worker started")
    
    // Simulate some work
    time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
    
    logger.WithFields(map[string]interface{}{
        "worker_id": id,
        "duration":  rand.Intn(100),
    }).Info("Worker completed")
}

func main() {
    logger := howl.New(howl.DefaultConfig())
    
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go worker(i, logger, &wg)
    }
    
    wg.Wait()
    logger.Info("All workers completed")
}
```

## Examples

Check the `examples` directory for more detailed examples:

- `examples/basic/basic.go`: Basic usage with different log levels and formats
- `examples/context/context_example.go`: Context integration example

## License

Howl is released under the MIT License. See the [LICENSE](LICENSE) file for details.

