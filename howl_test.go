package howl

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestLogLevels(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.level.String())
		}
	}
}

func TestBasicLogging(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(DefaultConfig())
	logger.SetOutput(buf)

	logger.Info("Test message")

	output := buf.String()
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected log to contain 'Test message', got: %s", output)
	}
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected log to contain 'INFO', got: %s", output)
	}
}

func TestWithField(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(DefaultConfig())
	logger.SetOutput(buf)

	logger.WithField("key", "value").Info("Test message")

	output := buf.String()
	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected log to contain 'key=value', got: %s", output)
	}
}

func TestWithFields(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(DefaultConfig())
	logger.SetOutput(buf)

	logger.WithFields(map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}).Info("Test message")

	output := buf.String()
	if !strings.Contains(output, "key1=value1") {
		t.Errorf("Expected log to contain 'key1=value1', got: %s", output)
	}
	if !strings.Contains(output, "key2=42") {
		t.Errorf("Expected log to contain 'key2=42', got: %s", output)
	}
}

func TestWithError(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(DefaultConfig())
	logger.SetOutput(buf)

	err := &testError{"test error"}
	logger.WithError(err).Error("Error occurred")

	output := buf.String()
	if !strings.Contains(output, "error=test error") {
		t.Errorf("Expected log to contain 'error=test error', got: %s", output)
	}
}

func TestJSONLogging(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(DefaultConfig().WithJSON(true).WithColor(false))
	logger.SetOutput(buf)

	logger.WithField("key", "value").Info("Test message")

	output := buf.String()

	// Check if the output contains JSON
	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Fatalf("Output doesn't contain JSON: %s", output)
	}

	// Check if it contains our field and message
	if !strings.Contains(output, "\"message\":") {
		t.Errorf("JSON doesn't contain message field: %s", output)
	}
	if !strings.Contains(output, "\"key\":") {
		t.Errorf("JSON doesn't contain key field: %s", output)
	}
	if !strings.Contains(output, "\"value\"") {
		t.Errorf("JSON doesn't contain value: %s", output)
	}
}

func TestContextLogging(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(DefaultConfig())
	logger.SetOutput(buf)

	// Add a context extractor
	logger = logger.WithContextExtractor(func(ctx context.Context) map[string]interface{} {
		fields := make(map[string]interface{})
		if requestID, ok := ctx.Value("request_id").(string); ok {
			fields["request_id"] = requestID
		}
		return fields
	})

	// Create a context with a request ID
	ctx := context.WithValue(context.Background(), "request_id", "123456")

	// Log with context
	logger.InfoContext(ctx, "Test with context")

	output := buf.String()
	if !strings.Contains(output, "request_id=123456") {
		t.Errorf("Expected log to contain 'request_id=123456', got: %s", output)
	}
}

func TestContextWithJSONLogging(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(DefaultConfig().WithJSON(true).WithColor(false))
	logger.SetOutput(buf)

	// Add a context extractor
	logger = logger.WithContextExtractor(func(ctx context.Context) map[string]interface{} {
		fields := make(map[string]interface{})
		if requestID, ok := ctx.Value("request_id").(string); ok {
			fields["request_id"] = requestID
		}
		return fields
	})

	// Create a context with a request ID
	ctx := context.WithValue(context.Background(), "request_id", "123456")

	// Log with context
	logger.InfoContext(ctx, "Test with context")

	output := buf.String()

	// Check if the output contains JSON
	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Fatalf("Output doesn't contain JSON: %s", output)
	}

	// Check if it contains our context field and message
	if !strings.Contains(output, "\"message\":") {
		t.Errorf("JSON doesn't contain message field: %s", output)
	}
	if !strings.Contains(output, "\"request_id\":") {
		t.Errorf("JSON doesn't contain request_id field: %s", output)
	}
	if !strings.Contains(output, "\"123456\"") {
		t.Errorf("JSON doesn't contain request ID value: %s", output)
	}
}

func TestMultipleContextExtractors(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(DefaultConfig())
	logger.SetOutput(buf)

	// Add multiple context extractors
	logger = logger.WithContextExtractor(func(ctx context.Context) map[string]interface{} {
		fields := make(map[string]interface{})
		if requestID, ok := ctx.Value("request_id").(string); ok {
			fields["request_id"] = requestID
		}
		return fields
	})

	logger = logger.WithContextExtractor(func(ctx context.Context) map[string]interface{} {
		fields := make(map[string]interface{})
		if user, ok := ctx.Value("user").(string); ok {
			fields["user"] = user
		}
		return fields
	})

	// Create a context with both request ID and user
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "123456")
	ctx = context.WithValue(ctx, "user", "admin")

	// Log with context
	logger.InfoContext(ctx, "Test with multiple extractors")

	output := buf.String()
	if !strings.Contains(output, "request_id=123456") {
		t.Errorf("Expected log to contain 'request_id=123456', got: %s", output)
	}
	if !strings.Contains(output, "user=admin") {
		t.Errorf("Expected log to contain 'user=admin', got: %s", output)
	}
}

// Helper types for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
