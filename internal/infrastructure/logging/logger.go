// Package logging provides a pluggable logger interface for the Huntress client library.
package logging

import (
	"context"
	"io"
	"log"
	"os"
)

// Level represents the log level for a logger.
type Level int

const (
	// LevelDebug enables debug-level logging.
	LevelDebug Level = iota
	// LevelInfo enables info-level logging.
	LevelInfo
	// LevelWarn enables warning-level logging.
	LevelWarn
	// LevelError enables error-level logging.
	LevelError
	// LevelFatal enables fatal-level logging.
	LevelFatal
)

// Field represents a structured log field.
type Field struct {
	Key   string
	Value interface{}
}

// Logger is the interface for structured logging.
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithContext(ctx context.Context) Logger
	WithFields(fields ...Field) Logger
}

// NoopLogger is a logger that does nothing (default).
type NoopLogger struct{}

// Debug does nothing for NoopLogger.
// Debug does nothing for NoopLogger.
func (NoopLogger) Debug(msg string, fields ...Field) {
	// intentionally left blank (noop)
	_ = msg
	_ = fields
}

// Info does nothing for NoopLogger.
// Info does nothing for NoopLogger.
func (NoopLogger) Info(msg string, fields ...Field) {
	// intentionally left blank (noop)
	_ = msg
	_ = fields
}

// Warn does nothing for NoopLogger.
// Warn does nothing for NoopLogger.
func (NoopLogger) Warn(msg string, fields ...Field) {
	// intentionally left blank (noop)
	_ = msg
	_ = fields
}

// Error does nothing for NoopLogger.
// Error does nothing for NoopLogger.
func (NoopLogger) Error(msg string, fields ...Field) {
	// intentionally left blank (noop)
	_ = msg
	_ = fields
}

// Fatal does nothing for NoopLogger.
// Fatal does nothing for NoopLogger.
func (NoopLogger) Fatal(msg string, fields ...Field) {
	// intentionally left blank (noop)
	_ = msg
	_ = fields
}

// WithContext returns NoopLogger for context.
func (NoopLogger) WithContext(context.Context) Logger { return NoopLogger{} }

// WithFields returns NoopLogger for fields.
func (NoopLogger) WithFields(...Field) Logger { return NoopLogger{} }

var (
	globalLogger Logger = NoopLogger{}
)

// SetLogger sets the global logger for the library.
func SetLogger(l Logger) {
	if l == nil {
		globalLogger = NoopLogger{}
	} else {
		globalLogger = l
	}
}

// L returns the global logger.
func L() Logger {
	return globalLogger
}

// StandardLogger is a simple implementation using the standard library log package.
type StandardLogger struct {
	level  Level
	logger *log.Logger
	fields []Field
}

// New creates a new StandardLogger with the given level, writer, and fields.
func New(level Level, w io.Writer, fields []Field) *StandardLogger {
	if w == nil {
		w = os.Stdout
	}
	return &StandardLogger{
		level:  level,
		logger: log.New(w, "", log.LstdFlags),
		fields: fields,
	}
}

func (l *StandardLogger) logf(lvl Level, prefix, msg string, fields ...Field) {
	if lvl < l.level {
		return
	}
	allFields := append(l.fields, fields...)
	l.logger.Printf("[%s] %s %v", prefix, msg, allFields)
}

// Debug logs a debug-level message.
func (l *StandardLogger) Debug(msg string, fields ...Field) {
	l.logf(LevelDebug, "DEBUG", msg, fields...)
}

// Info logs an info-level message.
func (l *StandardLogger) Info(msg string, fields ...Field) { l.logf(LevelInfo, "INFO", msg, fields...) }

// Warn logs a warning-level message.
func (l *StandardLogger) Warn(msg string, fields ...Field) { l.logf(LevelWarn, "WARN", msg, fields...) }

// Error logs an error-level message.
func (l *StandardLogger) Error(msg string, fields ...Field) {
	l.logf(LevelError, "ERROR", msg, fields...)
}

// Fatal logs a fatal-level message and exits the program.
func (l *StandardLogger) Fatal(msg string, fields ...Field) {
	l.logf(LevelFatal, "FATAL", msg, fields...)
	os.Exit(1)
}

// WithContext returns a logger with context (no-op for StandardLogger).
func (l *StandardLogger) WithContext(_ context.Context) Logger { return l }

// WithFields returns a logger with additional fields.
func (l *StandardLogger) WithFields(fields ...Field) Logger {
	return &StandardLogger{
		level:  l.level,
		logger: l.logger,
		fields: append(l.fields, fields...),
	}
}

// String creates a string log field.
func String(key, val string) Field { return Field{Key: key, Value: val} }

// Int creates an int log field.
func Int(key string, val int) Field { return Field{Key: key, Value: val} }

// Bool creates a bool log field.
func Bool(key string, val bool) Field { return Field{Key: key, Value: val} }

// Error creates an error log field.
func Error(key string, err error) Field { return Field{Key: key, Value: err} }
