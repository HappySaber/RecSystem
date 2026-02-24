package adapter

import (
	"api-gateway/internal/events/producer"
	"api-gateway/internal/events/schemas"
	"context"
)

type UserActionProducer struct {
	kafka *producer.Producer
	topic string
}

func NewUserActionProducer(kafka *producer.Producer, topic string) *UserActionProducer {
	return &UserActionProducer{
		kafka: kafka,
		topic: topic,
	}
}

func (p *UserActionProducer) Publish(
	ctx context.Context,
	topic, key string,
	event schemas.UserActionEvent,
) error {
	return p.kafka.Publish(ctx, p.topic, key, event)
}
