package recommendation_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"rec-system-microservice/internal/services/recommendation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockProvider struct{ mock.Mock }

func (m *mockProvider) GetUserRecommendations(ctx context.Context, userID string, limit int) ([]string, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockProvider) GetRecommendationsByGenres(ctx context.Context, genres []string, limit int) ([]string, error) {
	args := m.Called(ctx, genres, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockProvider) GetSimilarContent(ctx context.Context, contentID string, limit int) ([]string, error) {
	args := m.Called(ctx, contentID, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockProvider) GetTrendingContent(ctx context.Context, limit int) ([]string, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockProvider) GetPopularContent(ctx context.Context, limit int) ([]string, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]string), args.Error(1)
}

func newTestRecommendation() (*recommendation.Recommendation, *mockProvider) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	provider := new(mockProvider)
	svc := recommendation.New(log, provider, time.Hour)
	return svc, provider
}

func TestGetRecommendations_Success(t *testing.T) {
	svc, provider := newTestRecommendation()

	expected := []string{"content-1", "content-2"}
	provider.On("GetUserRecommendations", mock.Anything, "user-1", 10).
		Return(expected, nil).Once()

	result, err := svc.GetRecommendations(context.Background(), "user-1", 10)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
	provider.AssertExpectations(t)
}

func TestGetRecommendations_Error(t *testing.T) {
	svc, provider := newTestRecommendation()

	provider.On("GetUserRecommendations", mock.Anything, "user-1", 10).
		Return([]string{}, errors.New("redis unavailable")).Once()

	_, err := svc.GetRecommendations(context.Background(), "user-1", 10)

	require.Error(t, err)
	assert.ErrorContains(t, err, "redis unavailable")
	provider.AssertExpectations(t)
}

func TestGetRecommendationsByGenres_Success(t *testing.T) {
	svc, provider := newTestRecommendation()

	genres := []string{"Action", "Drama"}
	expected := []string{"content-1"}
	provider.On("GetRecommendationsByGenres", mock.Anything, genres, 5).
		Return(expected, nil).Once()

	result, err := svc.GetRecommendationsByGenres(context.Background(), genres, 5)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
	provider.AssertExpectations(t)
}

func TestGetRecommendationsByGenres_Error(t *testing.T) {
	svc, provider := newTestRecommendation()

	provider.On("GetRecommendationsByGenres", mock.Anything, mock.Anything, 5).
		Return([]string{}, errors.New("db error")).Once()

	_, err := svc.GetRecommendationsByGenres(context.Background(), []string{"Action"}, 5)

	require.Error(t, err)
	provider.AssertExpectations(t)
}

func TestGetSimilarContent_Success(t *testing.T) {
	svc, provider := newTestRecommendation()

	expected := []string{"content-2", "content-3"}
	provider.On("GetSimilarContent", mock.Anything, "content-1", 5).
		Return(expected, nil).Once()

	result, err := svc.GetSimilarContent(context.Background(), "content-1", 5)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	provider.AssertExpectations(t)
}

func TestGetTrendingContent_Success(t *testing.T) {
	svc, provider := newTestRecommendation()

	provider.On("GetTrendingContent", mock.Anything, 10).
		Return([]string{"content-1", "content-2"}, nil).Once()

	result, err := svc.GetTrendingContent(context.Background(), 10)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	provider.AssertExpectations(t)
}

func TestGetPopularContent_Success(t *testing.T) {
	svc, provider := newTestRecommendation()

	provider.On("GetPopularContent", mock.Anything, 10).
		Return([]string{"content-1"}, nil).Once()

	result, err := svc.GetPopularContent(context.Background(), 10)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	provider.AssertExpectations(t)
}
