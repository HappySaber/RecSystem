package consumer

import (
	"log/slog"
	"time"

	"notifications/internal/config"
)

// TestConfig возвращает конфиг с несуществующим брокером для тестов
func TestConfig() config.KafkaConfig {
	return config.KafkaConfig{
		Brokers:     []string{"localhost:9999"},
		Topic:       "test-topic",
		GroupID:     "test-group",
		ReadTimeout: 100 * time.Millisecond,
	}
}

// NewForTest создаёт consumer для тестов
func NewForTest(cfg config.KafkaConfig, log *slog.Logger) *KafkaConsumer {
	return New(cfg, log)
}
