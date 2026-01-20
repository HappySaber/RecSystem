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
}

func NewAIRecommendation(
	log *slog.Logger,
	client AIClient,
	prompts PromptBuilder,
	parser AIResponseParser,
) *AIRecommendation {
	return &AIRecommendation{log, client, prompts, parser}
}

type AIRecommendationProvider interface {
	GetImplicitRecommendations(
		ctx context.Context,
		userID string,
		limit int,
	) ([]string, error)

	GetExplicitRecommendations(
		ctx context.Context,
		query string,
		limit int,
	) ([]string, error)
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
		return nil, err
	}

	return a.parser.ParseContentIDs(raw)
}
