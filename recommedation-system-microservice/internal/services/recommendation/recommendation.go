package recommendation

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Recommendation struct {
	log      *slog.Logger
	engine   RecommendationProvider
	tokenTTL time.Duration
}

func New(
	log *slog.Logger,
	engine RecommendationProvider,
	tokenTTL time.Duration,
) *Recommendation {
	return &Recommendation{
		log:      log,
		engine:   engine,
		tokenTTL: tokenTTL,
	}
}

type RecommendationProvider interface {
	GetUserRecommendations(ctx context.Context, userID string, limit int) ([]string, error)
	GetRecommendationsByGenres(ctx context.Context, genres []string, limit int) ([]string, error)
	GetSimilarContent(ctx context.Context, contentID string, limit int) ([]string, error)
	GetTrendingContent(ctx context.Context, limit int) ([]string, error)
	GetPopularContent(ctx context.Context, limit int) ([]string, error)
}

func (r *Recommendation) GetRecommendations(
	ctx context.Context,
	userID string,
	limit int,
) ([]string, error) {

	const op = "recommendation.GetRecommendations"

	log := r.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("getting user recommendations")

	recs, err := r.engine.GetUserRecommendations(ctx, userID, limit)
	if err != nil {
		log.Error("failed to get recommendations", "error", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("recommendations retrieved", "count", len(recs))

	return recs, nil
}

func (r *Recommendation) GetRecommendationsByGenres(
	ctx context.Context,
	genres []string,
	limit int,
) ([]string, error) {

	const op = "recommendation.GetRecommendationsByGenres"

	log := r.log.With(
		slog.String("op", op),
		slog.Any("genres", genres),
	)

	log.Info("getting recommendations by genres")

	recs, err := r.engine.GetRecommendationsByGenres(ctx, genres, limit)
	if err != nil {
		log.Error("failed to get recommendations by genres", "error", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("genre recommendations retrieved", "count", len(recs))

	return recs, nil
}

func (r *Recommendation) GetSimilarContent(
	ctx context.Context,
	contentID string,
	limit int,
) ([]string, error) {

	const op = "recommendation.GetSimilarContent"

	log := r.log.With(
		slog.String("op", op),
		slog.String("content_id", contentID),
	)

	log.Info("getting similar content")

	recs, err := r.engine.GetSimilarContent(ctx, contentID, limit)
	if err != nil {
		log.Error("failed to get similar content", "error", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("similar content retrieved", "count", len(recs))

	return recs, nil
}

func (r *Recommendation) GetTrendingContent(
	ctx context.Context,
	limit int,
) ([]string, error) {

	const op = "recommendation.GetTrendingContent"

	log := r.log.With(
		slog.String("op", op),
	)

	log.Info("getting trending content")

	recs, err := r.engine.GetTrendingContent(ctx, limit)
	if err != nil {
		log.Error("failed to get trending content", "error", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("trending content retrieved", "count", len(recs))

	return recs, nil
}

func (r *Recommendation) GetPopularContent(
	ctx context.Context,
	limit int,
) ([]string, error) {

	const op = "recommendation.GetPopularContent"

	log := r.log.With(
		slog.String("op", op),
	)

	log.Info("getting popular content")

	recs, err := r.engine.GetPopularContent(ctx, limit)
	if err != nil {
		log.Error("failed to get popular content", "error", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("popular content retrieved", "count", len(recs))

	return recs, nil
}
