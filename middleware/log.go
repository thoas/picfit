package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoas/picfit/constants"
	"log/slog"
	"time"
)

func NewLogger(logger *slog.Logger) gin.HandlerFunc {
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
			slog.Duration("latency", latency),
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
	}
}
