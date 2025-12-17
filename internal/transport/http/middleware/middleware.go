package middleware

import (
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestLog := log.With(
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_addr", c.Request.RemoteAddr),
		)

		requestLog.Info("Request started")

		c.Set("logger", requestLog)
		c.Next()
		requestLog.Info("Request completed")
	}
}
