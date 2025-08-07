package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.SugaredLogger {
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		logLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "timestamp",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := config.Build()
	return logger.Sugar()
}

func GinLogger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Infow("Request handled",
				"status", c.Writer.Status(),
				"method", c.Request.Method,
				"path", path,
				"query", query,
				"ip", c.ClientIP(),
				"latency", latency.String(),
				"user-agent", c.Request.UserAgent(),
			)
		}
	}
}
