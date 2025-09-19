package core

import (
	"fmt"
	"log/slog"
	"strings"
)

// Level is a logger level.
type Level int8

// LevelKey is logger level key.
const LevelKey = "level"

const (
	// LevelTrace is logger trace level (most verbose).
	LevelTrace Level = iota - 2
	// LevelDebug is logger debug level.
	LevelDebug
	// LevelInfo is logger info level.
	LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn
	// LevelError is logger error level.
	LevelError
	// LevelFatal is logger fatal level.
	LevelFatal
)

// Key returns the string representation of the log level.
func (l Level) Key() string {
	return LevelKey
}

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Logger is a logger interface.
type Logger interface {
	Log(level Level, keyvals ...any) error
	With(keyvals ...any) Logger
}

// StdLogger is a standard logger implementation.
type StdLogger struct {
	logger *slog.Logger
}

// NewStdLogger creates a new standard logger.
func NewStdLogger(writer *slog.Logger) Logger {
	return &StdLogger{logger: writer}
}

// Log logs a message with the given level and key-value pairs.
func (s *StdLogger) Log(level Level, keyvals ...any) error {
	// Simple implementation for now
	if len(keyvals) == 0 {
		return nil
	}

	// Build log message
	var msg strings.Builder
	msg.WriteString(level.String())
	msg.WriteString(" ")

	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			msg.WriteString(fmt.Sprint(keyvals[i]))
			msg.WriteString("=")
			msg.WriteString(fmt.Sprint(keyvals[i+1]))
			msg.WriteString(" ")
		}
	}

	//nolint:noctx // no context needed
	s.logger.Info(msg.String())
	return nil
}

// With returns a new logger with the given key-value pairs.
func (s *StdLogger) With(keyvals ...any) Logger {
	return &logger{
		logger: s,
		prefix: keyvals,
	}
}

type logger struct {
	logger Logger
	prefix []any
}

func (c *logger) Log(level Level, keyvals ...any) error {
	kvs := make([]any, 0, len(c.prefix)+len(keyvals))
	kvs = append(kvs, c.prefix...)
	kvs = append(kvs, keyvals...)
	return c.logger.Log(level, kvs...)
}

func (c *logger) With(keyvals ...any) Logger {
	newPrefix := make([]any, 0, len(c.prefix)+len(keyvals))
	newPrefix = append(newPrefix, c.prefix...)
	newPrefix = append(newPrefix, keyvals...)
	return &logger{
		logger: c.logger,
		prefix: newPrefix,
	}
}

// NoOpLogger is a no-operation logger that does nothing.
type NoOpLogger struct{}

// Log logs a message with the given level and key-value pairs (no-op implementation).
func (n *NoOpLogger) Log(_ Level, _ ...any) error {
	return nil
}

// With returns a new logger with the given key-value pairs (no-op implementation).
func (n *NoOpLogger) With(_ ...any) Logger {
	return n
}
