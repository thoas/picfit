package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/constants"
	"log/slog"
	"runtime"
	"time"
)

func NewLogger(cfg *config.Config, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		requestID := uuid.New().String()
		c.Header("X-Request-ID", requestID)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constants.RequestIDCtx, requestID))

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		attributes := []slog.Attr{
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.String("latency", latency.String()),
			slog.String("user-agent", c.Request.UserAgent()),
			slog.Time("time", end.UTC()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.LogAttrs(c.Request.Context(), slog.LevelError, e, attributes...)
			}
		} else {
			logger.LogAttrs(c.Request.Context(), slog.LevelInfo, path, attributes...)
		}
		if cfg.Debug {
			logMemStats(c.Request.Context(), logger)

		}
	}
}

func logMemStats(ctx context.Context, logger *slog.Logger) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	attributes := []slog.Attr{
		slog.String("alloc", fmt.Sprintf("%v MiB", bToMb(m.Alloc))),
		slog.String("totalAlloc", fmt.Sprintf("%v MiB", bToMb(m.TotalAlloc))),
		slog.String("sys", fmt.Sprintf("%v MiB", bToMb(m.Sys))),
		slog.String("numGC", fmt.Sprintf("%v", m.NumGC)),
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "memory stats", attributes...)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
