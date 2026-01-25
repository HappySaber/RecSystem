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
	return &AIRecommendation{
		log:     log,
		client:  client,
		prompts: prompts,
		parser:  parser,
	}
}

type AIClient interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type PromptBuilder interface {
	BuildExplicit(query string, limit int) string
}

type AIResponseParser interface {
	ParseContentIDs(raw string) ([]string, error)
}

// type AIRecommendationProvider interface {
// 	GetImplicitRecommendations(
// 		ctx context.Context,
// 		userID string,
// 		limit int,
// 	) ([]string, error)

// 	GetExplicitRecommendations(
// 		ctx context.Context,
// 		query string,
// 		limit int,
// 	) ([]string, error)
// }

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
	return nil, errors.New("not implemented")
}
