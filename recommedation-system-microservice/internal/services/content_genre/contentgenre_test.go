package contentgenre_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	contentgenre "rec-system-microservice/internal/services/content_genre"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSaver struct{ mock.Mock }

func (m *mockSaver) SaveContentGenres(contentID string, genres []string) error {
	args := m.Called(contentID, genres)
	return args.Error(0)
}

func newTestContentGenre() (*contentgenre.ContentGenre, *mockSaver) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	saver := new(mockSaver)
	svc := contentgenre.New(log, saver)
	return svc, saver
}

func TestSaveContentGenres_Success(t *testing.T) {
	svc, saver := newTestContentGenre()

	genres := []string{"Action", "Drama", "Thriller"}
	saver.On("SaveContentGenres", "content-1", genres).
		Return(nil).Once()

	err := svc.SaveContentGenres(context.Background(), "content-1", genres)

	require.NoError(t, err)
	saver.AssertExpectations(t)
}

func TestSaveContentGenres_EmptyGenres(t *testing.T) {
	svc, saver := newTestContentGenre()

	// пустой список жанров — всё равно сохраняем, валидация на уровне выше
	saver.On("SaveContentGenres", "content-1", []string{}).
		Return(nil).Once()

	err := svc.SaveContentGenres(context.Background(), "content-1", []string{})

	require.NoError(t, err)
	saver.AssertExpectations(t)
}

func TestSaveContentGenres_Error(t *testing.T) {
	svc, saver := newTestContentGenre()

	saver.On("SaveContentGenres", "content-1", mock.Anything).
		Return(errors.New("redis unavailable")).Once()

	err := svc.SaveContentGenres(context.Background(), "content-1", []string{"Action"})

	require.Error(t, err)
	assert.ErrorContains(t, err, "redis unavailable")
	saver.AssertExpectations(t)
}
