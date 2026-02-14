package producer

import (
	"context"
	"encoding/json"
	"sso-microservice/internal/domain/events"

	"github.com/segmentio/kafka-go"
)

func (p *KafkaProducer) SendUserRegistered(
	ctx context.Context,
	event events.UserRegisteredEvent,
) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Email),
		Value: data,
	})
}
