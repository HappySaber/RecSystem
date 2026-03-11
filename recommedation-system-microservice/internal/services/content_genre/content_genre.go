package contentgenre

import (
	"context"
	"log/slog"
)

type ContentGenre struct {
	log               *slog.Logger
	contentGenreSaver ContentGenreSaver
}

type ContentGenreSaver interface {
	SaveContentGenres(contentID string, genres []string) error
}

func New(
	log *slog.Logger,
	contentGenreSaver ContentGenreSaver,
) *ContentGenre {
	return &ContentGenre{
		log:               log,
		contentGenreSaver: contentGenreSaver,
	}
}

func (cg *ContentGenre) SaveContentGenres(
	ctx context.Context,
	contentID string,
	genres []string,
) error {
	const op = "ContentGenre.SaveContentGenre"

	log := cg.log.With(
		slog.String("op", op),
	)
	log.Info("saving content genres", slog.String("content_id", contentID))

	err := cg.contentGenreSaver.SaveContentGenres(contentID, genres)
	if err != nil {
		log.Error("failed to save content genres", slog.String("error", err.Error()))
		return err
	}
	return nil
}
