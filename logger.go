package gologix

import (
	"log/slog"
)

// LoggerInterface defines the minimum logging interface required for basic logging.
// This interface provides the core logging methods but does not include contextual
// logging capabilities. When used for connections, additional connection information
// will not be included in log messages unless LoggerInterfaceWith is used.
type LoggerInterface interface {
	// Debug logs a message at debug level with optional key-value pairs.
	Debug(msg string, args ...any)

	// Error logs a message at error level with optional key-value pairs.
	Error(msg string, args ...any)

	// Warn logs a message at warning level with optional key-value pairs.
	Warn(msg string, args ...any)

	// Info logs a message at info level with optional key-value pairs.
	Info(msg string, args ...any)
}

// LoggerInterfaceWith extends LoggerInterface with contextual logging capabilities.
// When used for connections, this interface allows controllerIp details to be
// automatically included in all log messages.
type LoggerInterfaceWith interface {
	LoggerInterface

	// With returns a new logger with the given key-value pairs added to the context.
	// This is useful for adding connection-specific information to all subsequent logs.
	//
	// Used in connect.go
	With(args ...any) LoggerInterfaceWith
}

// Logger implements both LoggerInterface and LoggerInterfaceWith.
// It wraps the standard library's slog.Logger to provide the implemented interfaces.
type Logger struct {
	internalLogger *slog.Logger
}

func NewLogger() LoggerInterface {
	return &Logger{
		internalLogger: slog.Default(),
	}
}

func (l *Logger) SetLogger(logger *slog.Logger) {
	l.internalLogger = logger
}

func (l *Logger) With(args ...any) LoggerInterfaceWith {
	return &Logger{
		internalLogger: l.internalLogger.With(args...),
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.internalLogger.Debug(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.internalLogger.Error(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.internalLogger.Warn(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.internalLogger.Info(msg, args...)
}
