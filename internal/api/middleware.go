package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.Int("size", c.Writer.Size()),
		}

		ctx := c.Request.Context()

		if status >= 500 {
			s.log.WithCtx(ctx).Error("http server error", fields...)
		} else if status >= 400 {
			s.log.WithCtx(ctx).Warn("http client error", fields...)
		} else {
			s.log.WithCtx(ctx).Info("http request", fields...)
		}
	}
}
