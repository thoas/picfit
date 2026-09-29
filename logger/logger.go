package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type LogHandler struct {
	level       slog.Leveler
	sloghandler slog.Handler
	contextKeys []string
}

func (h LogHandler) Handle(ctx context.Context, record slog.Record) error {
	if !h.Enabled(ctx, record.Level) {
		return nil
	}
	contextfields := getContextFields(ctx, h.contextKeys)
	if len(contextfields) > 0 {
		record.AddAttrs(contextfields...)
	}
	return h.sloghandler.Handle(ctx, record)
}

func (h LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return LogHandler{
		level:       h.level,
		sloghandler: h.sloghandler.WithAttrs(attrs),
		contextKeys: h.contextKeys,
	}
}

func (h LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{
		sloghandler: h.sloghandler.WithGroup(name),
	}
}

func (h LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func getContextFields(ctx context.Context, keys []string) []slog.Attr {
	var fields []slog.Attr
	for _, k := range keys {
		value := ctx.Value(k)
		if value != nil {
			fields = append(fields, slog.Any(k, value))
		}
	}
	return fields
}

func New(cfg Config) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	var opts = slog.HandlerOptions{Level: level}
	return slog.New(LogHandler{
		contextKeys: cfg.ContextKeys,
		level:       opts.Level,
		sloghandler: slog.NewJSONHandler(os.Stderr, &opts),
	})

}
