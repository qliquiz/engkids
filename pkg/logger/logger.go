package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewLogger(appName string) (*logrus.Logger, error) {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()

	logger.SetOutput(io.MultiWriter(os.Stdout, file))

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	logger.SetLevel(logrus.InfoLevel)

	return logger, nil
}

func LoggingMiddleware(log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		entry := log.WithFields(logrus.Fields{
			"method":     c.Method(),
			"path":       c.OriginalURL(),
			"status":     c.Response().StatusCode(),
			"latency_ms": latency.Milliseconds(),
			"ip":         c.IP(),
			"request_id": c.Locals("requestid"),
			"user_agent": c.Get("User-Agent"),
			"time":       time.Now().Format(time.RFC3339Nano),
		})

		if err != nil {
			entry.WithField("error", err.Error()).Error("request failed")
		} else {
			entry.Info("request completed")
		}

		return err
	}
}
