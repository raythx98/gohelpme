package logger

import (
	"context"
	"fmt"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"log/slog"
	"os"
)

var (
	defaultLogger *Default
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	defaultLogger = &Default{}
}

// Default is a struct that wraps the default logger from the "log" package.
// It satisfied logger.ILogger interface and provides additional methods for logging with context and options.
//
// It is used as a fallback logger when the zap logger is not available.
type Default struct{}

// NewDefault creates a new instance of Default logger using Go's native logger.
func NewDefault() *Default {
	return defaultLogger
}

func (d *Default) GetInstance() interface{} {
	return slog.Default()
}

func (d *Default) Debug(ctx context.Context, msg string, options ...Field) {
	slog.Debug(msg, fmt.Sprintf("%+v", reqctx.GetValue(ctx)), GetMapFromFields(options...))
}

func (d *Default) Info(ctx context.Context, msg string, options ...Field) {
	slog.Info(msg, fmt.Sprintf("%+v", reqctx.GetValue(ctx)), GetMapFromFields(options...))
}

func (d *Default) Warn(ctx context.Context, msg string, options ...Field) {
	slog.Warn(msg, fmt.Sprintf("%+v", reqctx.GetValue(ctx)), GetMapFromFields(options...))
}

func (d *Default) Error(ctx context.Context, msg string, options ...Field) {
	slog.Error(msg, fmt.Sprintf("%+v", reqctx.GetValue(ctx)), GetMapFromFields(options...))
}

func (d *Default) Fatal(ctx context.Context, msg string, options ...Field) {
	slog.Error(msg, fmt.Sprintf("%+v", reqctx.GetValue(ctx)), GetMapFromFields(options...))
	// TODO: Should we exit here?
	os.Exit(1)
}

func (d *Default) Panic(ctx context.Context, msg string, options ...Field) {
	slog.Error(fmt.Sprintf("%+v", reqctx.GetValue(ctx)), msg, GetMapFromFields(options...))
	// TODO: Should we exit here?
	os.Exit(1)
}
