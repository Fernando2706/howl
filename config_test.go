package howl

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Check default values
	if config.Level != InfoLevel {
		t.Errorf("Expected default level to be InfoLevel, got %v", config.Level)
	}
	if config.JSON != false {
		t.Errorf("Expected default JSON to be false, got %v", config.JSON)
	}
	if config.Color != true {
		t.Errorf("Expected default Color to be true, got %v", config.Color)
	}
	if config.Timestamp != true {
		t.Errorf("Expected default Timestamp to be true, got %v", config.Timestamp)
	}
}

func TestWithLevel(t *testing.T) {
	config := DefaultConfig().WithLevel(DebugLevel)
	if config.Level != DebugLevel {
		t.Errorf("Expected level to be DebugLevel, got %v", config.Level)
	}
}

func TestWithJSON(t *testing.T) {
	config := DefaultConfig().WithJSON(true)
	if config.JSON != true {
		t.Errorf("Expected JSON to be true, got %v", config.JSON)
	}
}

func TestWithColor(t *testing.T) {
	config := DefaultConfig().WithColor(false)
	if config.Color != false {
		t.Errorf("Expected Color to be false, got %v", config.Color)
	}
}

func TestWithTimestamp(t *testing.T) {
	config := DefaultConfig().WithTimestamp(false)
	if config.Timestamp != false {
		t.Errorf("Expected Timestamp to be false, got %v", config.Timestamp)
	}
}

func TestChainedConfig(t *testing.T) {
	config := DefaultConfig().
		WithLevel(ErrorLevel).
		WithJSON(true).
		WithColor(false).
		WithTimestamp(false)

	// Check all values
	if config.Level != ErrorLevel {
		t.Errorf("Expected level to be ErrorLevel, got %v", config.Level)
	}
	if config.JSON != true {
		t.Errorf("Expected JSON to be true, got %v", config.JSON)
	}
	if config.Color != false {
		t.Errorf("Expected Color to be false, got %v", config.Color)
	}
	if config.Timestamp != false {
		t.Errorf("Expected Timestamp to be false, got %v", config.Timestamp)
	}
}
