package consumer

import (
	"context"
	"encoding/json"
	"rec-system-microservice/internal/schemas"
	contentgenre "rec-system-microservice/internal/services/content_genre"
)

func ContentGenreHandler(service *contentgenre.ContentGenre) MessageHandler {
	return func(ctx context.Context, msg []byte) error {
		var event schemas.ContentGenre

		if err := json.Unmarshal(msg, &event); err != nil {
			return err
		}

		return service.SaveContentGenre(ctx, event.ContentID, event.Genres)
	}
}
