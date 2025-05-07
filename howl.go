package howl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Logger represents the main logging entity
type Logger struct {
	mu            sync.Mutex
	level         Level
	config        Config
	writer        io.Writer
	fields        map[string]interface{}
	ctxExtractors []ContextExtractor
}

// ContextExtractor is a function that extracts fields from a context
type ContextExtractor func(ctx context.Context) map[string]interface{}

// New creates a new logger with the provided configuration
func New(config Config) *Logger {
	return &Logger{
		level:         config.Level,
		config:        config,
		writer:        os.Stdout,
		fields:        make(map[string]interface{}),
		ctxExtractors: make([]ContextExtractor, 0),
	}
}

// Default creates a new logger with default configuration
func Default() *Logger {
	return New(DefaultConfig())
}

// SetOutput sets the output destination for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer = w
}

// SetLevel sets the minimum level for logging
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// WithField returns a new logger with the field added to the logging context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := &Logger{
		level:  l.level,
		config: l.config,
		writer: l.writer,
		fields: make(map[string]interface{}, len(l.fields)+1),
	}

	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	newLogger.fields[key] = value

	return newLogger
}

// WithFields returns a new logger with the fields added to the logging context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := &Logger{
		level:  l.level,
		config: l.config,
		writer: l.writer,
		fields: make(map[string]interface{}, len(l.fields)+len(fields)),
	}

	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// WithError returns a new logger with the error added to the logging context
func (l *Logger) WithError(err error) *Logger {
	return l.WithField("error", err.Error())
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DebugLevel, nil, msg, args...)
}

// Info logs a message at info level
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(InfoLevel, nil, msg, args...)
}

// Warn logs a message at warning level
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WarnLevel, nil, msg, args...)
}

// Error logs a message at error level
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ErrorLevel, nil, msg, args...)
}

// Fatal logs a message at fatal level and then exits with status code 1
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(FatalLevel, nil, msg, args...)
	os.Exit(1)
}

// DebugContext logs a message at debug level with context
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	l.log(DebugLevel, ctx, msg, args...)
}

// InfoContext logs a message at info level with context
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	l.log(InfoLevel, ctx, msg, args...)
}

// WarnContext logs a message at warning level with context
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	l.log(WarnLevel, ctx, msg, args...)
}

// ErrorContext logs a message at error level with context
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	l.log(ErrorLevel, ctx, msg, args...)
}

// FatalContext logs a message at fatal level with context and then exits with status code 1
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...interface{}) {
	l.log(FatalLevel, ctx, msg, args...)
	os.Exit(1)
}

// WithContextExtractor adds a context extractor to the logger
func (l *Logger) WithContextExtractor(extractor ContextExtractor) *Logger {
	newLogger := l.clone()
	newLogger.ctxExtractors = append(newLogger.ctxExtractors, extractor)
	return newLogger
}

// clone creates a copy of the logger with the same handlers and formatter
func (l *Logger) clone() *Logger {
	fields := make(map[string]interface{}, len(l.fields))
	for k, v := range l.fields {
		fields[k] = v
	}

	extractors := make([]ContextExtractor, len(l.ctxExtractors))
	copy(extractors, l.ctxExtractors)

	return &Logger{
		level:         l.level,
		config:        l.config,
		writer:        l.writer,
		fields:        fields,
		ctxExtractors: extractors,
	}
}

// log handles the actual logging
func (l *Logger) log(level Level, ctx context.Context, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	isFile := false
	if file, ok := l.writer.(*os.File); ok {
		info, err := file.Stat()
		if err == nil && !info.Mode().IsRegular() {
			isFile = false
		} else {
			isFile = true
		}
	}

	config := l.config
	if isFile {
		config.Color = false
	}

	var output string

	if config.Timestamp {
		timestamp := time.Now().Format(time.RFC3339)
		output += timestamp + " "
	}

	contextFields := make(map[string]interface{})
	if ctx != nil {
		for _, extractor := range l.ctxExtractors {
			fields := extractor(ctx)
			for k, v := range fields {
				contextFields[k] = v
			}
		}
	}

	levelStr := level.String()
	if config.Color {
		levelColor := getLevelColor(level)
		output += "[" + levelColor + levelStr + "\033[0m] "
	} else {
		output += "[" + levelStr + "] "
	}

	if !config.JSON {
		output += msg
	}

	allFields := make(map[string]interface{}, len(l.fields)+len(contextFields))
	for k, v := range l.fields {
		allFields[k] = v
	}
	for k, v := range contextFields {
		allFields[k] = v
	}

	if len(allFields) > 0 {
		if config.JSON {
			jsonData := make(map[string]interface{})

			jsonData["message"] = msg

			for k, v := range allFields {
				jsonData[k] = v
			}

			jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
			if err != nil {
				output += fmt.Sprintf(" Error formatting JSON: %s", err.Error())
			} else {
				jsonStr := string(jsonBytes)
				if config.Color {
					coloredJSON := colorizeJSON(jsonStr)
					output += "\n" + coloredJSON
				} else {
					output += "\n" + string(jsonBytes)
				}
			}
		} else {
			for k, v := range allFields {
				output += fmt.Sprintf(" %s=%v", k, v)
			}
		}
	}

	fmt.Fprintln(l.writer, output)

	if level == FatalLevel {
		os.Exit(1)
	}
}

// getLevelColor returns the ANSI color code for the given level
func getLevelColor(level Level) string {
	switch level {
	case DebugLevel:
		return "\033[36m"
	case InfoLevel:
		return "\033[32m"
	case WarnLevel:
		return "\033[33m"
	case ErrorLevel:
		return "\033[31m"
	case FatalLevel:
		return "\033[35m"
	default:
		return "\033[0m"
	}
}

// ANSI color codes for JSON formatting
const (
	keyColor   = "\033[36m"
	strColor   = "\033[32m"
	numColor   = "\033[33m"
	boolColor  = "\033[35m"
	nullColor  = "\033[37m"
	braceColor = "\033[34m"
	commaColor = "\033[37m"
	resetColor = "\033[0m"
)

// colorizeJSON adds ANSI color codes to a JSON string to make it more readable
func colorizeJSON(jsonStr string) string {
	var result strings.Builder
	inString := false
	escapeNext := false

	keyColor := "\033[36m"         // Cyan for keys
	stringColor := "\033[32m"      // Green for string values
	numberColor := "\033[33m"      // Yellow for numeric values
	boolColor := "\033[35m"        // Magenta for boolean values
	nullColor := "\033[37m"        // White for null values
	syntaxColor := "\033[34m"      // Blue for braces and brackets
	punctuationColor := "\033[37m" // White for commas and colons
	resetColor := "\033[0m"        // Reset color

	var tokenBuffer strings.Builder
	collectingToken := false
	isKey := false

	for i, char := range jsonStr {
		if inString {
			if escapeNext {
				escapeNext = false
				result.WriteString(string(char))
				continue
			}

			if char == '\\' {
				escapeNext = true
				result.WriteString(string(char))
				continue
			}

			if char == '"' {
				inString = false
				result.WriteString(string(char) + resetColor)
				continue
			}

			result.WriteString(string(char))
			continue
		}

		switch char {
		case '"':
			isKey = false
			// Check if this is a key (followed by a colon)
			for j := i + 1; j < len(jsonStr); j++ {
				if jsonStr[j] == '"' {
					// Found the end of the string
					for k := j + 1; k < len(jsonStr); k++ {
						if jsonStr[k] == ' ' || jsonStr[k] == '\n' || jsonStr[k] == '\t' || jsonStr[k] == '\r' {
							continue
						}
						if jsonStr[k] == ':' {
							isKey = true
						}
						break
					}
					break
				}
			}

			inString = true
			if isKey {
				result.WriteString(keyColor + string(char))
			} else {
				result.WriteString(stringColor + string(char))
			}

		case '{', '[':
			result.WriteString(syntaxColor + string(char) + resetColor)

		case '}', ']':
			result.WriteString(syntaxColor + string(char) + resetColor)

		case ':':
			result.WriteString(punctuationColor + string(char) + resetColor)

		case ',':
			result.WriteString(punctuationColor + string(char) + resetColor)

		case ' ', '\n', '\t', '\r':
			result.WriteString(string(char))

		default:
			if !collectingToken {
				collectingToken = true
				tokenBuffer.Reset()
			}
			tokenBuffer.WriteRune(char)

			token := tokenBuffer.String()
			nextChar := byte(0)
			if i+1 < len(jsonStr) {
				nextChar = jsonStr[i+1]
			}

			if nextChar == ',' || nextChar == '}' || nextChar == ']' || nextChar == ' ' || nextChar == '\n' || nextChar == '\t' || nextChar == '\r' {
				if token == "true" || token == "false" {
					result.WriteString(boolColor + token + resetColor)
				} else if token == "null" {
					result.WriteString(nullColor + token + resetColor)
				} else {
					result.WriteString(numberColor + token + resetColor)
				}
				collectingToken = false
				continue
			}

			result.WriteString(string(char))
			collectingToken = false
		}
	}

	return result.String()
}
