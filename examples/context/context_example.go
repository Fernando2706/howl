package main

import (
	"context"
	"fmt"
	"time"

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

// UserExtractor extracts user information from the context
func UserExtractor(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})
	if user, ok := ctx.Value("user").(string); ok {
		fields["user"] = user
	}
	return fields
}

// simulateHTTPRequest simulates an HTTP request with context
func simulateHTTPRequest(logger *howl.Logger, requestID string, user string) {
	// Create a context with request ID and user
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", requestID)
	ctx = context.WithValue(ctx, "user", user)

	// Log the start of the request
	logger.InfoContext(ctx, "Request started")

	// Simulate some processing
	time.Sleep(100 * time.Millisecond)

	// Log with additional fields
	logger.WithField("path", "/api/users").InfoContext(ctx, "Processing request")

	// Simulate more processing
	time.Sleep(50 * time.Millisecond)

	// Log the completion of the request
	logger.WithFields(map[string]interface{}{
		"status":   200,
		"duration": 150,
	}).InfoContext(ctx, "Request completed")
}

func main() {
	// Create a logger with context extractors
	logger := howl.New(
		howl.DefaultConfig().
			WithJSON(true).
			WithColor(true).
			WithTimestamp(true),
	)

	// Add context extractors
	logger = logger.WithContextExtractor(RequestIDExtractor)
	logger = logger.WithContextExtractor(UserExtractor)

	fmt.Println("Starting context example...")

	// Simulate multiple HTTP requests
	simulateHTTPRequest(logger, "req-123", "admin")
	simulateHTTPRequest(logger, "req-456", "user")
	simulateHTTPRequest(logger, "req-789", "guest")

	// Create a logger with text format
	textLogger := howl.New(
		howl.DefaultConfig().
			WithJSON(false).
			WithColor(true),
	)
	textLogger = textLogger.WithContextExtractor(RequestIDExtractor)
	textLogger = textLogger.WithContextExtractor(UserExtractor)

	fmt.Println("\nText format with context:")
	simulateHTTPRequest(textLogger, "req-abc", "support")

	fmt.Println("\nContext example completed")
}
