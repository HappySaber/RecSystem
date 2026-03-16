package suite

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// InsertContent вставляет запись в таблицу content и возвращает сгенерированный UUID
func (st *Suite) InsertContent(t *testing.T, contentType, externalSource, externalID, title string) string {
	t.Helper()

	var id string
	err := st.DB.QueryRow(`
        INSERT INTO content (type, external_source, external_id, title, description, poster_url, release_date)
        VALUES ($1, $2, $3, $4, 'test description', 'http://poster.url', '2024-01-01')
        RETURNING id`,
		contentType, externalSource, externalID, title,
	).Scan(&id)
	require.NoError(t, err, "failed to insert content")

	t.Cleanup(func() {
		st.DB.Exec(`DELETE FROM content WHERE id = $1`, id)
	})

	return id
}

// InsertMovieDetails вставляет запись в movies_details
func (st *Suite) InsertMovieDetails(t *testing.T, contentID string, tmdbID int) {
	t.Helper()

	_, err := st.DB.Exec(`
		INSERT INTO movies_details (content_id, tmdb_id, raw_data)
		VALUES ($1, $2, '{}')`,
		contentID, tmdbID,
	)
	require.NoError(t, err, "failed to insert movie details")
}

// InsertAnimeDetails вставляет запись в anime_details
func (st *Suite) InsertAnimeDetails(t *testing.T, contentID string, anilistID int) {
	t.Helper()

	_, err := st.DB.Exec(`
		INSERT INTO anime_details (content_id, anilist_id, raw_data)
		VALUES ($1, $2, '{}')`,
		contentID, anilistID,
	)
	require.NoError(t, err, "failed to insert anime details")
}

// InsertGameDetails вставляет запись в games_details
func (st *Suite) InsertGameDetails(t *testing.T, contentID string, igdbID int) {
	t.Helper()

	_, err := st.DB.Exec(`
		INSERT INTO games_details (content_id, igdb_id, raw_data)
		VALUES ($1, $2, '{}')`,
		contentID, igdbID,
	)
	require.NoError(t, err, "failed to insert game details")
}

// InsertSeriesDetails вставляет запись в series_details
func (st *Suite) InsertSeriesDetails(t *testing.T, contentID string, tmdbID int) {
	t.Helper()

	_, err := st.DB.Exec(`
		INSERT INTO series_details (content_id, tmdb_id, raw_data)
		VALUES ($1, $2, '{}')`,
		contentID, tmdbID,
	)
	require.NoError(t, err, "failed to insert series details")
}

// InsertBookDetails вставляет запись в books_details
func (st *Suite) InsertBookDetails(t *testing.T, contentID string) {
	t.Helper()

	_, err := st.DB.Exec(`
		INSERT INTO books_details (content_id, raw_data)
		VALUES ($1, '{}')`,
		contentID,
	)
	require.NoError(t, err, "failed to insert book details")
}
