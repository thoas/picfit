package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
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

func LogMemStats(ctx context.Context, msg string, logger *slog.Logger) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	attributes := []slog.Attr{
		slog.String("alloc", fmt.Sprintf("%v MiB", bToMb(m.Alloc))),
		slog.String("heap-alloc", fmt.Sprintf("%v MiB", bToMb(m.HeapAlloc))),
		slog.String("total-alloc", fmt.Sprintf("%v MiB", bToMb(m.TotalAlloc))),
		slog.String("sys", fmt.Sprintf("%v MiB", bToMb(m.Sys))),
		slog.Int("numgc", int(m.NumGC)),
		slog.Int("total-goroutine", runtime.NumGoroutine()),
	}
	logger.LogAttrs(ctx, slog.LevelInfo, msg, attributes...)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
