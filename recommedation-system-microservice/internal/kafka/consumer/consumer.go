package consumer

import (
	"context"
	"encoding/json"
	"rec-system-microservice/internal/domain/models"
	useractions "rec-system-microservice/internal/services/user_actions"

	"github.com/segmentio/kafka-go"
)

type UserActionConsumer struct {
	reader  *kafka.Reader
	service *useractions.UserActionTracker
}

func NewUserActionConsumer(r *kafka.Reader, s *useractions.UserActionTracker) *UserActionConsumer {
	return &UserActionConsumer{
		reader:  r,
		service: s,
	}
}

func (c *UserActionConsumer) Start(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var event models.UserActionEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			continue
		}

		if err := c.service.TrackUserAction(ctx, event); err != nil {
			// логируем, но не падаем
			continue
		}
	}
}
