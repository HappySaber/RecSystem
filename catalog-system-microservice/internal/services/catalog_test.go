package catalog_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"catalog-microservice/internal/domain/models"
	catalog "catalog-microservice/internal/services"
	"catalog-microservice/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// мок провайдера — реализует интерфейс CatalogProvider
type mockProvider struct {
	mock.Mock
}

func (m *mockProvider) GetContent(ctx context.Context, id string) (models.Content, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Content), args.Error(1)
}

func (m *mockProvider) GetContentByIDs(ctx context.Context, ids []string) ([]models.ContentShort, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]models.ContentShort), args.Error(1)
}

func (m *mockProvider) FindContentByExternal(ctx context.Context, externalID, externalSource string) (models.Content, error) {
	args := m.Called(ctx, externalID, externalSource)
	return args.Get(0).(models.Content), args.Error(1)
}

func (m *mockProvider) AnimeDetails(ctx context.Context, id string) (models.AnimeDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.AnimeDetails), args.Error(1)
}

func (m *mockProvider) AllAnimeDetails(ctx context.Context) ([]models.AnimeDetails, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.AnimeDetails), args.Error(1)
}

func (m *mockProvider) MovieDetails(ctx context.Context, id string) (models.MovieDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.MovieDetails), args.Error(1)
}

func (m *mockProvider) AllMovieDetails(ctx context.Context) ([]models.MovieDetails, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.MovieDetails), args.Error(1)
}

func (m *mockProvider) SeriesDetails(ctx context.Context, id string) (models.SeriesDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.SeriesDetails), args.Error(1)
}

func (m *mockProvider) AllSeriesDetails(ctx context.Context) ([]models.SeriesDetails, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.SeriesDetails), args.Error(1)
}

func (m *mockProvider) BookDetails(ctx context.Context, id string) (models.BookDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.BookDetails), args.Error(1)
}

func (m *mockProvider) AllBookDetails(ctx context.Context) ([]models.BookDetails, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.BookDetails), args.Error(1)
}

func (m *mockProvider) GameDetails(ctx context.Context, id string) (models.GameDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.GameDetails), args.Error(1)
}

func (m *mockProvider) AllGameDetails(ctx context.Context) ([]models.GameDetails, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.GameDetails), args.Error(1)
}

// хелпер создания сервиса с моком
func newTestService() (*catalog.Catalog, *mockProvider) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	provider := new(mockProvider)
	svc := catalog.New(log, provider)
	return svc, provider
}

// GetContent

func TestGetContent_Success(t *testing.T) {
	svc, provider := newTestService()

	expected := models.Content{
		ID:             "abc-123",
		Title:          "Inception",
		ExternalSource: "tmdb",
		ExternalID:     "27205",
	}

	provider.On("GetContent", mock.Anything, "abc-123").
		Return(expected, nil).Once()

	result, err := svc.GetContent(context.Background(), "abc-123")

	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Title, result.Title)
	provider.AssertExpectations(t)
}

func TestGetContent_NotFound(t *testing.T) {
	svc, provider := newTestService()

	provider.On("GetContent", mock.Anything, "missing-id").
		Return(models.Content{}, storage.ErrNotFound).Once()

	_, err := svc.GetContent(context.Background(), "missing-id")

	require.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNotFound)
	provider.AssertExpectations(t)
}

func TestGetContent_InternalError(t *testing.T) {
	svc, provider := newTestService()

	provider.On("GetContent", mock.Anything, "abc-123").
		Return(models.Content{}, errors.New("db connection lost")).Once()

	_, err := svc.GetContent(context.Background(), "abc-123")

	require.Error(t, err)
	// не ErrNotFound — значит внутренняя ошибка
	assert.False(t, errors.Is(err, storage.ErrNotFound))
	provider.AssertExpectations(t)
}

// GetContentByIDs

func TestGetContentByIDs_Success(t *testing.T) {
	svc, provider := newTestService()

	expected := []models.ContentShort{
		{ID: "id-1", Type: "movie", Title: "Inception"},
		{ID: "id-2", Type: "anime", Title: "Naruto"},
	}

	provider.On("GetContentByIDs", mock.Anything, []string{"id-1", "id-2"}).
		Return(expected, nil).Once()

	result, err := svc.GetContentByIDs(context.Background(), []string{"id-1", "id-2"})

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Inception", result[0].Title)
	provider.AssertExpectations(t)
}

func TestGetContentByIDs_NotFound(t *testing.T) {
	svc, provider := newTestService()

	provider.On("GetContentByIDs", mock.Anything, []string{"missing"}).
		Return([]models.ContentShort{}, storage.ErrNotFound).Once()

	_, err := svc.GetContentByIDs(context.Background(), []string{"missing"})

	require.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNotFound)
	provider.AssertExpectations(t)
}

// FindContentByExternal

func TestFindContentByExternal_Success(t *testing.T) {
	svc, provider := newTestService()

	expected := models.Content{
		ID:             "abc-123",
		Title:          "The Matrix",
		ExternalSource: "tmdb",
		ExternalID:     "603",
	}

	provider.On("FindContentByExternal", mock.Anything, "603", "tmdb").
		Return(expected, nil).Once()

	result, err := svc.FindContentByExternal(context.Background(), "603", "tmdb")

	require.NoError(t, err)
	assert.Equal(t, "The Matrix", result.Title)
	assert.Equal(t, "tmdb", result.ExternalSource)
	provider.AssertExpectations(t)
}

func TestFindContentByExternal_NotFound(t *testing.T) {
	svc, provider := newTestService()

	provider.On("FindContentByExternal", mock.Anything, "nonexistent", "tmdb").
		Return(models.Content{}, storage.ErrNotFound).Once()

	_, err := svc.FindContentByExternal(context.Background(), "nonexistent", "tmdb")

	require.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNotFound)
	provider.AssertExpectations(t)
}

// GetMovieDetails

func TestGetMovieDetails_Success(t *testing.T) {
	svc, provider := newTestService()

	expected := models.MovieDetails{
		ContentID: "movie-1",
		TmdbID:    550,
	}

	provider.On("MovieDetails", mock.Anything, "movie-1").
		Return(expected, nil).Once()

	result, err := svc.GetMovieDetails(context.Background(), "movie-1")

	require.NoError(t, err)
	assert.Equal(t, 550, result.TmdbID)
	assert.Equal(t, "movie-1", result.ContentID)
	provider.AssertExpectations(t)
}

func TestGetMovieDetails_NotFound(t *testing.T) {
	svc, provider := newTestService()

	provider.On("MovieDetails", mock.Anything, "missing").
		Return(models.MovieDetails{}, storage.ErrNotFound).Once()

	_, err := svc.GetMovieDetails(context.Background(), "missing")

	require.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNotFound)
	provider.AssertExpectations(t)
}

// GetAllMovieDetails

func TestGetAllMovieDetails_Success(t *testing.T) {
	svc, provider := newTestService()

	expected := []models.MovieDetails{
		{ContentID: "m1", TmdbID: 550},
		{ContentID: "m2", TmdbID: 551},
	}

	provider.On("AllMovieDetails", mock.Anything).
		Return(expected, nil).Once()

	result, err := svc.GetAllMovieDetails(context.Background())

	require.NoError(t, err)
	assert.Len(t, result, 2)
	provider.AssertExpectations(t)
}

func TestGetAllMovieDetails_Empty(t *testing.T) {
	svc, provider := newTestService()

	provider.On("AllMovieDetails", mock.Anything).
		Return([]models.MovieDetails{}, nil).Once()

	result, err := svc.GetAllMovieDetails(context.Background())

	require.NoError(t, err)
	assert.Empty(t, result)
	provider.AssertExpectations(t)
}

// GetAnimeDetails

func TestGetAnimeDetails_Success(t *testing.T) {
	svc, provider := newTestService()

	expected := models.AnimeDetails{
		ContentID: "anime-1",
		AniListID: 20,
	}

	provider.On("AnimeDetails", mock.Anything, "anime-1").
		Return(expected, nil).Once()

	result, err := svc.GetAnimeDetails(context.Background(), "anime-1")

	require.NoError(t, err)
	assert.Equal(t, 20, result.AniListID)
	provider.AssertExpectations(t)
}

func TestGetAnimeDetails_NotFound(t *testing.T) {
	svc, provider := newTestService()

	provider.On("AnimeDetails", mock.Anything, "missing").
		Return(models.AnimeDetails{}, storage.ErrNotFound).Once()

	_, err := svc.GetAnimeDetails(context.Background(), "missing")

	require.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNotFound)
	provider.AssertExpectations(t)
}

// GetAllAnimeDetails

func TestGetAllAnimeDetails_Success(t *testing.T) {
	svc, provider := newTestService()

	expected := []models.AnimeDetails{
		{ContentID: "a1", AniListID: 20},
		{ContentID: "a2", AniListID: 21},
	}

	provider.On("AllAnimeDetails", mock.Anything).
		Return(expected, nil).Once()

	result, err := svc.GetAllAnimeDetails(context.Background())

	require.NoError(t, err)
	assert.Len(t, result, 2)
	provider.AssertExpectations(t)
}

// GetSeriesDetails

func TestGetSeriesDetails_Success(t *testing.T) {
	svc, provider := newTestService()

	expected := models.SeriesDetails{
		ContentID: "series-1",
		TmdbID:    1396,
	}

	provider.On("SeriesDetails", mock.Anything, "series-1").
		Return(expected, nil).Once()

	result, err := svc.GetSeriesDetails(context.Background(), "series-1")

	require.NoError(t, err)
	assert.Equal(t, 1396, result.TmdbID)
	provider.AssertExpectations(t)
}

func TestGetSeriesDetails_NotFound(t *testing.T) {
	svc, provider := newTestService()

	provider.On("SeriesDetails", mock.Anything, "missing").
		Return(models.SeriesDetails{}, storage.ErrNotFound).Once()

	_, err := svc.GetSeriesDetails(context.Background(), "missing")

	require.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNotFound)
	provider.AssertExpectations(t)
}

// GetAllSeriesDetails

func TestGetAllSeriesDetails_Success(t *testing.T) {
	svc, provider := newTestService()

	expected := []models.SeriesDetails{
		{ContentID: "s1", TmdbID: 1396},
		{ContentID: "s2", TmdbID: 1399},
	}

	provider.On("AllSeriesDetails", mock.Anything).
		Return(expected, nil).Once()

	result, err := svc.GetAllSeriesDetails(context.Background())

	require.NoError(t, err)
	assert.Len(t, result, 2)
	provider.AssertExpectations(t)
}

func TestGetAllSeriesDetails_Empty(t *testing.T) {
	svc, provider := newTestService()

	provider.On("AllSeriesDetails", mock.Anything).
		Return([]models.SeriesDetails{}, nil).Once()

	result, err := svc.GetAllSeriesDetails(context.Background())

	require.NoError(t, err)
	assert.Empty(t, result)
	provider.AssertExpectations(t)
}
