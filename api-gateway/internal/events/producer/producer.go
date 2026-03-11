package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		}}
}

func (p *Producer) Publish(ctx context.Context, topic, key string, event any) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	fmt.Println("KAFKA SEND:", topic, string(payload))
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	})
}
