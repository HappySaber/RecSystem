package consumer

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, msg []byte) error

type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

func NewConsumer(r *kafka.Reader, h MessageHandler) Consumer {
	return Consumer{
		reader:  r,
		handler: h,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("kafka read error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if err := c.handler(ctx, msg.Value); err != nil {
			log.Println("handler error:", err)
		}
	}
}
