package howl

// Config holds the configuration options for the logger
type Config struct {
	Level     Level
	JSON      bool
	Color     bool
	Timestamp bool
}

// DefaultConfig returns a default configuration for the logger
func DefaultConfig() Config {
	return Config{
		Level:     InfoLevel,
		JSON:      false,
		Color:     true,
		Timestamp: true,
	}
}

// WithLevel sets the log level and returns the updated config
func (c Config) WithLevel(level Level) Config {
	c.Level = level
	return c
}

// WithJSON enables or disables JSON formatting
func (c Config) WithJSON(enabled bool) Config {
	c.JSON = enabled
	return c
}

// WithColor enables or disables colored output
func (c Config) WithColor(enabled bool) Config {
	c.Color = enabled
	return c
}

// WithTimestamp enables or disables timestamps in log entries
func (c Config) WithTimestamp(enabled bool) Config {
	c.Timestamp = enabled
	return c
}
