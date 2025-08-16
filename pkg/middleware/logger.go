package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoggerMiddleware struct {
	Logger *zap.Logger
}

func NewLoggerMiddleware(logger *zap.Logger) LoggerMiddleware {
	return LoggerMiddleware{logger}
}

func (m *LoggerMiddleware) LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		m.Logger.Info("HTTP Request",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.Duration("duration", duration),
		)
	}
}
