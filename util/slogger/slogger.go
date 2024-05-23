package slogger

import (
	"context"
	"log/slog"
	"os"

	"github.com/raythx98/gohelpme/util/reqctx"
)

func Init() {
	handler := &ContextHandler{slog.NewJSONHandler(os.Stdout, nil)}
	slog.SetDefault(slog.New(handler))
}

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

type ContextHandler struct {
	slog.Handler
}

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	if reqCtx := reqctx.GetValue(ctx); reqCtx != nil {
		r.AddAttrs(slog.Any("context", reqctx.GetValue(ctx)))
	}

	return h.Handler.Handle(ctx, r)
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	v := make([]slog.Attr, 0)
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}

func GetLogLevel(err error) slog.Level {
	if err != nil {
		return slog.LevelError
	}
	return slog.LevelInfo
}
