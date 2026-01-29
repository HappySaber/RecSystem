package consumer

import (
	"context"
	"errors"
	"log/slog"

	"notifications/internal/config"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	cfg    config.KafkaConfig
	reader *kafka.Reader
	log    *slog.Logger
}

func New(cfg config.KafkaConfig, log *slog.Logger) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID,
		MaxWait:  cfg.ReadTimeout,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		cfg:    cfg,
		reader: reader,
		log:    log,
	}
}

func (kc *KafkaConsumer) Start(ctx context.Context, handler func(msg kafka.Message) error) {
	const op = "consumer.Start"

	log := kc.log.With(slog.String("op", op))

	go func() {
		defer func() {
			log.Info("closing kafka reader")
			_ = kc.reader.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				log.Info("context cancelled, stopping consumer")
				return
			default:
				msg, err := kc.reader.ReadMessage(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}

					log.Error("failed to read message", slog.Any("error", err))
					continue
				}

				if err := handler(msg); err != nil {
					log.Error(
						"handler failed",
						slog.Any("error", err),
						slog.Int64("offset", msg.Offset),
					)
				}
			}
		}
	}()
}
