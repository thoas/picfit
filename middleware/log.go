package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/constants"
)

func NewLogger(cfg *config.Config, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx       = c.Request.Context()
			start     = time.Now()
			path      = c.Request.URL.Path
			requestID = uuid.New().String()
		)
		c.Header("X-Request-ID", requestID)
		c.Request = c.Request.WithContext(context.WithValue(ctx, constants.RequestIDCtx, requestID))

		ctx = c.Request.Context()

		c.Next()

		end := time.Now()
		attributes := []slog.Attr{
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.Any("params", c.Request.URL.Query()),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Duration("duration", time.Since(start)),
			slog.String("user-agent", c.Request.UserAgent()),
			slog.Time("time", end.UTC()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.LogAttrs(ctx, slog.LevelError, e, attributes...)
			}
		} else {
			logger.LogAttrs(ctx, slog.LevelInfo, path, attributes...)
		}

		logMemStats(ctx, logger)
	}
}

func logMemStats(ctx context.Context, logger *slog.Logger) {
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
	logger.LogAttrs(ctx, slog.LevelInfo, "Memory stats", attributes...)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
