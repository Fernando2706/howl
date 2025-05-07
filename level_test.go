package howl

import (
	"bytes"
	"strings"
	"testing"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
		{Level(999), "UNKNOWN"}, // Test unknown level
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Level.String() for %d: expected %s, got %s", test.level, test.expected, test.level.String())
		}
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		levelStr string
		expected Level
	}{
		{"debug", DebugLevel},
		{"DEBUG", DebugLevel},
		{"info", InfoLevel},
		{"INFO", InfoLevel},
		{"warn", WarnLevel},
		{"WARN", WarnLevel},
		{"error", ErrorLevel},
		{"ERROR", ErrorLevel},
		{"fatal", FatalLevel},
		{"FATAL", FatalLevel},
		// For invalid input, ParseLevel returns InfoLevel as default
		{"invalid", InfoLevel},
	}

	for _, test := range tests {
		level := ParseLevel(test.levelStr)
		if level != test.expected {
			t.Errorf("ParseLevel(%q): expected %d, got %d", test.levelStr, test.expected, level)
		}
	}
}

func TestLevelColorization(t *testing.T) {
	// Just test that we can log with different levels without crashing
	config := DefaultConfig().WithColor(true).WithLevel(DebugLevel)
	logger := New(config)

	// Test each level
	buf := new(bytes.Buffer)
	logger.SetOutput(buf)

	// Test debug level
	logger.Debug("Debug message")
	output := buf.String()
	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Log output doesn't contain DEBUG level: %s", output)
	}

	// Clear buffer
	buf.Reset()

	// Test info level
	logger.Info("Info message")
	output = buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("Log output doesn't contain INFO level: %s", output)
	}

	// Clear buffer
	buf.Reset()

	// Test warn level
	logger.Warn("Warn message")
	output = buf.String()
	if !strings.Contains(output, "WARN") {
		t.Errorf("Log output doesn't contain WARN level: %s", output)
	}

	// Clear buffer
	buf.Reset()

	// Test error level
	logger.Error("Error message")
	output = buf.String()
	if !strings.Contains(output, "ERROR") {
		t.Errorf("Log output doesn't contain ERROR level: %s", output)
	}
}
