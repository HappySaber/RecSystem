package producer

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	log    *slog.Logger
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (kp *KafkaProducer) SendMessage(ctx context.Context, key, value []byte) error {
	return kp.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}
