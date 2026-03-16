package suite

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// InsertGenre вставляет жанр и возвращает его id
func (st *Suite) InsertGenre(t *testing.T, name string) int {
	t.Helper()

	var id int
	err := st.DB.QueryRow(`
        INSERT INTO genres (name)
        VALUES ($1)
        ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
        RETURNING id`,
		name,
	).Scan(&id)
	require.NoError(t, err)

	t.Cleanup(func() {
		st.DB.Exec(`DELETE FROM genres WHERE id = $1`, id)
	})
	return id
}

// InsertContentGenre привязывает контент к жанру
func (st *Suite) InsertContentGenre(t *testing.T, contentID string, genreID int) {
	t.Helper()

	_, err := st.DB.Exec(`
        INSERT INTO content_genres (content_id, genre_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING`,
		contentID, genreID,
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		st.DB.Exec(`DELETE FROM content_genres WHERE content_id = $1`, contentID)
	})
}

// InsertUserAction вставляет действие пользователя
// actionCode — VIEW | LIKE | DISLIKE | RATE | FAVORITE
func (st *Suite) InsertUserAction(t *testing.T, userID, contentID, actionCode string) {
	t.Helper()

	_, err := st.DB.Exec(`
        INSERT INTO user_content_action (user_id, content_id, action_id)
        SELECT $1, $2, id FROM actions WHERE code = $3`,
		userID, contentID, actionCode,
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		st.DB.Exec(`
            DELETE FROM user_content_action
            WHERE user_id = $1 AND content_id = $2`,
			userID, contentID,
		)
	})
}
