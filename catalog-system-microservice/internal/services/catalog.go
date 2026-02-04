package catalog

import (
	"catalog-microservice/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type Catalog struct {
	log             *slog.Logger
	catalogProvider CatalogProvider
}

type CatalogProvider interface {
	GetContent(ctx context.Context, id string) (models.Content, error)
	GetContentByIDs(ctx context.Context, ids []string) ([]models.ContentShort, error)

	FindContentByExternal(ctx context.Context, externalID string) (models.Content, error)
	//AllAnimeDetails(ctx context.Context) ([]models.AnimeDetails, error)
	AnimeDetails(ctx context.Context, id string) (models.AnimeDetails, error)
	AllAnimeDetails(ctx context.Context) ([]models.AnimeDetails, error)
	BookDetails(ctx context.Context, id string) (models.BookDetails, error)
	AllBookDetails(ctx context.Context) ([]models.BookDetails, error)
	MovieDetails(ctx context.Context, id string) (models.MovieDetails, error)
	AllMovieDetails(ctx context.Context) ([]models.MovieDetails, error)
	SeriesDetails(ctx context.Context, id string) (models.SeriesDetails, error)
	AllSeriesDetails(ctx context.Context) ([]models.SeriesDetails, error)
	GameDetails(ctx context.Context, id string) (models.GameDetails, error)
	AllGameDetails(ctx context.Context) ([]models.GameDetails, error)
}

type AppProvider interface {
}

var (
	ErrInvalidID       = errors.New("invalid ID")
	ErrIDDoesNotExists = errors.New("ID does not exists")
)

func New(
	log *slog.Logger,
	catalogProvider CatalogProvider,
) *Catalog {
	return &Catalog{
		log:             log,
		catalogProvider: catalogProvider,
	}
}

func (c *Catalog) GetContent(ctx context.Context, id string) (models.Content, error) {
	const op = "catalog.GetContent"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	content, err := c.catalogProvider.GetContent(ctx, id)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return models.Content{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return models.Content{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Content founded")

	return content, nil
}

func (c *Catalog) GetContentByIDs(ctx context.Context, ids []string) ([]models.ContentShort, error) {
	const op = "catalog.GetContentByIDs"
	log := c.log.With(
		slog.String("op", op),
	)
	log.Info("trying to get contents by IDs")

	contents, err := c.catalogProvider.GetContentByIDs(ctx, ids)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("contents not found", "error", err.Error())
			return []models.ContentShort{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get contents", "error", err.Error())
		return []models.ContentShort{}, fmt.Errorf("%s: %w", op, err)
	}
	c.log.Info("Contents founded")
	return contents, nil
}

func (c *Catalog) FindContentByExternal(ctx context.Context, externalID string) (models.Content, error) {
	const op = "catalog.FindContentByExternal"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by external ID")

	content, err := c.catalogProvider.FindContentByExternal(ctx, externalID)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return models.Content{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return models.Content{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Content founded")
	return content, nil
}
func (c *Catalog) GetMovieDetails(ctx context.Context, id string) (models.MovieDetails, error) {
	const op = "catalog.GetMovieDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	movie, err := c.catalogProvider.MovieDetails(ctx, id)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return models.MovieDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return models.MovieDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Content founded")
	return movie, nil
}
func (c *Catalog) GetAnimeDetails(ctx context.Context, id string) (models.AnimeDetails, error) {
	const op = "catalog.GetAnimeDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	anime, err := c.catalogProvider.AnimeDetails(ctx, id)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return models.AnimeDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return models.AnimeDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Anime founded")
	return anime, nil
}
func (c *Catalog) GetGameDetails(ctx context.Context, id string) (models.GameDetails, error) {
	const op = "catalog.GetGameDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	game, err := c.catalogProvider.GameDetails(ctx, id)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return models.GameDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return models.GameDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Game founded")
	return game, nil
}
func (c *Catalog) GetSeriesDetails(ctx context.Context, id string) (models.SeriesDetails, error) {
	const op = "catalog.GetSeriesDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	series, err := c.catalogProvider.SeriesDetails(ctx, id)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return models.SeriesDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return models.SeriesDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Series founded")
	return series, nil
}
func (c *Catalog) GetBookDetails(ctx context.Context, id string) (models.BookDetails, error) {
	const op = "catalog.GetBookDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	book, err := c.catalogProvider.BookDetails(ctx, id)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return models.BookDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return models.BookDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("book founded")
	return book, nil
}
func (c *Catalog) GetAllMovieDetails(ctx context.Context) ([]models.MovieDetails, error) {
	const op = "catalog.GetAllMovieDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	movies, err := c.catalogProvider.AllMovieDetails(ctx)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return []models.MovieDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return []models.MovieDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("all movies founded")
	return movies, nil
}
func (c *Catalog) GetAllAnimeDetails(ctx context.Context) ([]models.AnimeDetails, error) {
	const op = "catalog.GetAllAnimeDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	animes, err := c.catalogProvider.AllAnimeDetails(ctx)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return []models.AnimeDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return []models.AnimeDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("all movies founded")
	return animes, nil
}
func (c *Catalog) GetAllGameDetails(ctx context.Context) ([]models.GameDetails, error) {
	const op = "catalog.GetAllGameDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	games, err := c.catalogProvider.AllGameDetails(ctx)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return []models.GameDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return []models.GameDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("all games founded")
	return games, nil
}
func (c *Catalog) GetAllSeriesDetails(ctx context.Context) ([]models.SeriesDetails, error) {
	const op = "catalog.GetAllSeriesDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	serieses, err := c.catalogProvider.AllSeriesDetails(ctx)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return []models.SeriesDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return []models.SeriesDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("all serieses founded")
	return serieses, nil
}
func (c *Catalog) GetAllBookDetails(ctx context.Context) ([]models.BookDetails, error) {
	const op = "catalog.GetAllBookDetails"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("trying to get content by ID")

	books, err := c.catalogProvider.AllBookDetails(ctx)
	if err != nil {
		if errors.Is(err, ErrIDDoesNotExists) {
			c.log.Warn("content not found", "error", err.Error())
			return []models.BookDetails{}, fmt.Errorf("%s: %w", op, ErrIDDoesNotExists)
		}
		c.log.Error("failed to get content", "error", err.Error())

		return []models.BookDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("all books founded")
	return books, nil
}
