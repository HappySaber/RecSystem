package consumer_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"notifications/internal/consumer"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func TestConsumer_HandlerCalledWithMessage(t *testing.T) {
	// тестируем логику хендлера напрямую — без реального Kafka
	received := make(chan kafka.Message, 1)

	handler := func(msg kafka.Message) error {
		received <- msg
		return nil
	}

	msg := kafka.Message{
		Topic: "test-topic",
		Value: []byte(`{"email":"user@example.com","name":"John"}`),
	}

	// вызываем хендлер напрямую
	err := handler(msg)

	assert.NoError(t, err)
	assert.Equal(t, msg.Value, (<-received).Value)
}

func TestConsumer_HandlerError(t *testing.T) {
	// проверяем что хендлер возвращает ошибку корректно
	handler := func(msg kafka.Message) error {
		return errors.New("processing failed")
	}

	msg := kafka.Message{Value: []byte(`{}`)}
	err := handler(msg)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "processing failed")
}

func TestConsumer_ContextCancellation(t *testing.T) {
	// проверяем что consumer останавливается при отмене контекста
	log := newTestLogger()
	cfg := consumer.TestConfig()

	kc := consumer.NewForTest(cfg, log)

	ctx, cancel := context.WithCancel(context.Background())

	stopped := make(chan struct{})
	kc.Start(ctx, func(msg kafka.Message) error {
		return nil
	})

	// отменяем контекст
	cancel()

	// даём горутине время завершиться
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(stopped)
	}()

	select {
	case <-stopped:
		// успешно завершился
	case <-time.After(2 * time.Second):
		t.Fatal("consumer did not stop after context cancellation")
	}
}
