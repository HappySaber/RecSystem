package producer

import (
	"catalog-microservice/internal/schemas"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

func (p *KafkaProducer) SendUserRegistered(
	ctx context.Context,
	event schemas.ContentGenre,
) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ContentID),
		Value: data,
	})
}
