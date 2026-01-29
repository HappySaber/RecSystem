package consumer

import (
	"context"
	"log/slog"
	"notifications/internal/config"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	config *config.Config
	reader *kafka.Reader
	log    *slog.Logger
}

func New(config *config.Config, topic, groupID string, brokers []string, log *slog.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		config: config,
		log:    log,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (kc *KafkaConsumer) Start(ctx context.Context, handler func(msg kafka.Message)) {
	const op = "consumer.Start"

	log := kc.log.With(
		slog.String("op", op),
	)

	go func() {
		defer kc.reader.Close()

		for {
			m, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				log.Error("failed to read message", "error", err.Error())
				break
			}
			handler(m)
		}
	}()
}
