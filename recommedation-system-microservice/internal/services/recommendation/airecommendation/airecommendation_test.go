package airecommendation_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"rec-system-microservice/internal/services/recommendation/airecommendation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAIClient struct{ mock.Mock }

func (m *mockAIClient) Complete(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

type mockPromptBuilder struct{ mock.Mock }

func (m *mockPromptBuilder) BuildExplicit(query string, limit int) string {
	args := m.Called(query, limit)
	return args.String(0)
}

func (m *mockPromptBuilder) BuildImplicit(genres map[string]float64, limit int) string {
	args := m.Called(genres, limit)
	return args.String(0)
}

type mockParser struct{ mock.Mock }

func (m *mockParser) ParseContentIDs(raw string) ([]string, error) {
	args := m.Called(raw)
	return args.Get(0).([]string), args.Error(1)
}

type mockGenreRepo struct{ mock.Mock }

func (m *mockGenreRepo) GetUserGenreWeights(ctx context.Context, userID string, limit int) (map[string]float64, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func newTestAI() (*airecommendation.AIRecommendation, *mockAIClient, *mockPromptBuilder, *mockParser, *mockGenreRepo) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := new(mockAIClient)
	prompts := new(mockPromptBuilder)
	parser := new(mockParser)
	repo := new(mockGenreRepo)
	svc := airecommendation.NewAIRecommendation(log, client, prompts, parser, repo)
	return svc, client, prompts, parser, repo
}

// GetExplicitRecommendations

func TestGetExplicitRecommendations_Success(t *testing.T) {
	svc, client, prompts, parser, _ := newTestAI()

	prompts.On("BuildExplicit", "anime action", 5).
		Return("some prompt").Once()

	client.On("Complete", mock.Anything, "some prompt").
		Return(`["Naruto","Bleach"]`, nil).Once()

	parser.On("ParseContentIDs", `["Naruto","Bleach"]`).
		Return([]string{"Naruto", "Bleach"}, nil).Once()

	result, err := svc.GetExplicitRecommendations(context.Background(), "anime action", 5)

	require.NoError(t, err)
	assert.Equal(t, []string{"Naruto", "Bleach"}, result)
	prompts.AssertExpectations(t)
	client.AssertExpectations(t)
	parser.AssertExpectations(t)
}

func TestGetExplicitRecommendations_AIClientError(t *testing.T) {
	svc, client, prompts, _, _ := newTestAI()

	prompts.On("BuildExplicit", mock.Anything, mock.Anything).
		Return("some prompt").Once()

	client.On("Complete", mock.Anything, "some prompt").
		Return("", errors.New("huggingface unavailable")).Once()

	_, err := svc.GetExplicitRecommendations(context.Background(), "anime action", 5)

	require.Error(t, err)
	assert.ErrorContains(t, err, "huggingface unavailable")
}

func TestGetExplicitRecommendations_ParseError(t *testing.T) {
	svc, client, prompts, parser, _ := newTestAI()

	prompts.On("BuildExplicit", mock.Anything, mock.Anything).
		Return("some prompt").Once()

	client.On("Complete", mock.Anything, "some prompt").
		Return("invalid json", nil).Once()

	parser.On("ParseContentIDs", "invalid json").
		Return([]string{}, airecommendation.ErrInvalidAIResponse).Once()

	_, err := svc.GetExplicitRecommendations(context.Background(), "anime action", 5)

	require.Error(t, err)
	assert.ErrorIs(t, err, airecommendation.ErrInvalidAIResponse)
}

// GetImplicitRecommendations

func TestGetImplicitRecommendations_Success(t *testing.T) {
	svc, client, prompts, parser, repo := newTestAI()

	genres := map[string]float64{"Action": 0.8, "Drama": 0.5}

	repo.On("GetUserGenreWeights", mock.Anything, "user-1", 5).
		Return(genres, nil).Once()

	prompts.On("BuildImplicit", genres, 10).
		Return("some prompt").Once()

	client.On("Complete", mock.Anything, "some prompt").
		Return(`["content-1","content-2"]`, nil).Once()

	parser.On("ParseContentIDs", `["content-1","content-2"]`).
		Return([]string{"content-1", "content-2"}, nil).Once()

	result, err := svc.GetImplicitRecommendations(context.Background(), "user-1", 10)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
	prompts.AssertExpectations(t)
	client.AssertExpectations(t)
	parser.AssertExpectations(t)
}

func TestGetImplicitRecommendations_NoGenres(t *testing.T) {
	svc, _, _, _, repo := newTestAI()

	// пустые жанры — возвращаем nil без вызова AI
	repo.On("GetUserGenreWeights", mock.Anything, "user-1", 5).
		Return(map[string]float64{}, nil).Once()

	result, err := svc.GetImplicitRecommendations(context.Background(), "user-1", 10)

	require.NoError(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetImplicitRecommendations_RepoError(t *testing.T) {
	svc, _, _, _, repo := newTestAI()

	repo.On("GetUserGenreWeights", mock.Anything, "user-1", 5).
		Return(map[string]float64{}, errors.New("redis unavailable")).Once()

	_, err := svc.GetImplicitRecommendations(context.Background(), "user-1", 10)

	require.Error(t, err)
	assert.ErrorContains(t, err, "redis unavailable")
}
