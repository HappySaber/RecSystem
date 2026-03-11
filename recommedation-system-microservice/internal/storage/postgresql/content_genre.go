package postgresql

import "fmt"

func (s *Storage) SaveContentGenres(contentID string, genres []string) error {
	const op = "postgresql.SaveContentGenres"

	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, genre := range genres {

		var genreID int

		err := tx.QueryRow(`
			INSERT INTO genres (name)
			VALUES ($1)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, genre).Scan(&genreID)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s insert genre: %w", op, err)
		}

		_, err = tx.Exec(`
			INSERT INTO content_genres (content_id, genre_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, contentID, genreID)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s insert content_genre: %w", op, err)
		}
	}

	return tx.Commit()
}
