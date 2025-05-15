package logger

import (
	"context"
	"os"
)

type contextKey struct{}

var loggerKey = &contextKey{}

func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}
	return New(os.Stdout) // 返回默认logger
}

func WithFields(ctx context.Context, fields Fields) Logger {
	return FromContext(ctx).WithFields(fields)
}
