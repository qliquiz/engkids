package logger

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Logger обертка над logrus с поддержкой отправки логов в Logstash
type Logger struct {
	*logrus.Logger
}

// LogstashHook хук для отправки логов в Logstash
type LogstashHook struct {
	conn    net.Conn
	appName string
}

// NewLogger создает новый экземпляр логгера
func NewLogger(appName string) (*Logger, error) {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	log.SetOutput(os.Stdout)

	// Получаем настройки из переменных окружения
	logstashHost := os.Getenv("LOGSTASH_HOST")
	logstashPort := os.Getenv("LOGSTASH_PORT")

	if logstashHost != "" && logstashPort != "" {
		hookAddress := fmt.Sprintf("%s:%s", logstashHost, logstashPort)
		hook, err := NewLogstashHook("tcp", hookAddress, appName)
		if err != nil {
			return nil, fmt.Errorf("failed to create logstash hook: %v", err)
		}
		log.AddHook(hook)
	}

	return &Logger{log}, nil
}

// NewLogstashHook создает новый хук для отправки логов в Logstash
func NewLogstashHook(network, address, appName string) (*LogstashHook, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &LogstashHook{
		conn:    conn,
		appName: appName,
	}, nil
}

// Fire отправляет запись лога в Logstash
func (hook *LogstashHook) Fire(entry *logrus.Entry) error {
	data := make(logrus.Fields, len(entry.Data)+5)

	for k, v := range entry.Data {
		data[k] = v
	}

	data["@timestamp"] = entry.Time.Format(time.RFC3339)
	data["message"] = entry.Message
	data["level"] = entry.Level.String()
	data["type"] = "engkids-log"
	data["app"] = hook.appName

	if hostname, err := os.Hostname(); err == nil {
		data["host"] = hostname
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal fields to JSON: %v", err)
	}

	_, err = hook.conn.Write(append(serialized, '\n'))
	return err
}

// Levels возвращает уровни логирования для хука
func (hook *LogstashHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// LoggingMiddleware создает middleware для логирования HTTP-запросов
func LoggingMiddleware(log *Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Вызов следующего обработчика
		err := c.Next()

		// Логируем результат запроса
		log.WithFields(logrus.Fields{
			"method":      c.Method(),
			"path":        c.Path(),
			"status":      c.Response().StatusCode(),
			"ip":          c.IP(),
			"user_agent":  c.Get("User-Agent"),
			"duration_ms": time.Since(start).Milliseconds(),
			"request_id":  c.GetRespHeader("X-Request-ID", ""),
		}).Info("HTTP request")

		return err
	}
}
