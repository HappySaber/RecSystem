package airecommendation

import (
	"context"
	"errors"
	"log/slog"
)

var (
	ErrInvalidAIResponse = errors.New("invalid AI response")
)

type AIRecommendation struct {
	log     *slog.Logger
	client  AIClient
	prompts PromptBuilder
	parser  AIResponseParser
	repo    GenreRepository
}

func NewAIRecommendation(
	log *slog.Logger,
	client AIClient,
	prompts PromptBuilder,
	parser AIResponseParser,
	repo GenreRepository,
) *AIRecommendation {
	return &AIRecommendation{
		log:     log,
		client:  client,
		prompts: prompts,
		parser:  parser,
		repo:    repo,
	}
}

type AIClient interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type PromptBuilder interface {
	BuildExplicit(query string, limit int) string
	BuildImplicit(genres map[string]float64, limit int) string
}

type AIResponseParser interface {
	ParseContentIDs(raw string) ([]string, error)
}

type GenreRepository interface {
	GetUserGenreWeights(
		ctx context.Context,
		userID string,
		limit int,
	) (map[string]float64, error)
}

func (a *AIRecommendation) GetExplicitRecommendations(
	ctx context.Context,
	query string,
	limit int,
) ([]string, error) {
	const op = "recommendation.GetExplicitRecommendations"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("getting explicit recommendations", slog.String("query", query))

	prompt := a.prompts.BuildExplicit(query, limit)

	log.Info("sending prompt to AI client", slog.String("prompt", prompt))
	raw, err := a.client.Complete(ctx, prompt)
	if err != nil {
		log.Error("failed to parse AI response",
			slog.String("raw", raw),
		)
		return nil, err
	}

	return a.parser.ParseContentIDs(raw)
}
func (a *AIRecommendation) GetImplicitRecommendations(
	ctx context.Context,
	userID string,
	limit int,
) ([]string, error) {

	const op = "recommendation.GetImplicitRecommendations"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting implicit recommendations", slog.String("userID", userID))

	genres, err := a.repo.GetUserGenreWeights(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	if len(genres) == 0 {
		log.Warn("user has no genre preferences")
		return nil, nil
	}

	prompt := a.prompts.BuildImplicit(genres, limit)

	log.Info("sending prompt to AI client", slog.String("prompt", prompt))

	raw, err := a.client.Complete(ctx, prompt)
	if err != nil {
		log.Error("failed to call AI client", slog.String("error", err.Error()))
		return nil, err
	}

	return a.parser.ParseContentIDs(raw)
}
