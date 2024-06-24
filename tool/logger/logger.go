package logger

import (
	"context"
)

//go:generate mockery --name ILogger
type ILogger interface {
	// Debug logs a message at DebugLevel.
	Debug(ctx context.Context, msg string, fields ...Field)
	// Info logs a message at InfoLevel.
	Info(ctx context.Context, msg string, fields ...Field)
	// Warn logs a message at WarnLevel.
	Warn(ctx context.Context, msg string, fields ...Field)
	// Error logs a message at ErrorLevel.
	Error(ctx context.Context, msg string, fields ...Field)
	// Fatal logs a message at FatalLevel.
	Fatal(ctx context.Context, msg string, fields ...Field)
	// Panic logs a message at PanicLevel.
	Panic(ctx context.Context, msg string, fields ...Field)
}
