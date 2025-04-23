package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoas/picfit/constants"
	loggerpkg "github.com/thoas/picfit/logger"
)

func NewLogger(logger *slog.Logger) gin.HandlerFunc {
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

		attributes := []slog.Attr{
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.Any("params", c.Request.URL.Query()),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Duration("duration", time.Since(start)),
			slog.String("user-agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.LogAttrs(ctx, slog.LevelError, e, attributes...)
			}
		} else {
			logger.LogAttrs(ctx, slog.LevelInfo, path, attributes...)
		}

		loggerpkg.WithMemStats(logger).Info("Memory stats")
	}
}
