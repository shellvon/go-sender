package core

import (
	"log"
	"strings"
)

// Level is a logger level.
type Level int8

// LevelKey is logger level key.
const LevelKey = "level"

const (
	// LevelDebug is logger debug level.
	LevelDebug Level = iota - 1
	// LevelInfo is logger info level.
	LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn
	// LevelError is logger error level.
	LevelError
	// LevelFatal is logger fatal level
	LevelFatal
)

func (l Level) Key() string {
	return LevelKey
}

func (l Level) String() string {
	switch l {
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

// StdLogger is a standard logger implementation
type StdLogger struct {
	logger *log.Logger
}

// NewStdLogger creates a new standard logger
func NewStdLogger(writer *log.Logger) Logger {
	return &StdLogger{logger: writer}
}

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
			msg.WriteString(keyvals[i].(string))
			msg.WriteString("=")
			msg.WriteString(keyvals[i+1].(string))
			msg.WriteString(" ")
		}
	}

	s.logger.Println(msg.String())
	return nil
}

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

// NoOpLogger is a no-operation logger that does nothing
type NoOpLogger struct{}

func (n *NoOpLogger) Log(level Level, keyvals ...any) error {
	return nil
}

func (n *NoOpLogger) With(keyvals ...any) Logger {
	return n
}
