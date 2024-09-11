package logger

import (
	"context"
	"log/slog"
)

type LogLevel int

func (l LogLevel) ToSlogLevel() slog.Level {
	return slog.Level(l)
}

type Attr = slog.Attr

type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)

	DebugContext(ctx context.Context, msg string, fields ...any)
	InfoContext(ctx context.Context, msg string, fields ...any)
	WarnContext(ctx context.Context, msg string, fields ...any)
	ErrorContext(ctx context.Context, msg string, fields ...any)

	Log(ctx context.Context, level LogLevel, msg string, fields ...any)

	String(key, value string) Attr

	With(args ...any) *slog.Logger
}
