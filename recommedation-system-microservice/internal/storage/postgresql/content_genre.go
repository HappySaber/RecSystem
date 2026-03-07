package postgresql

import "fmt"

func (s *Storage) SaveContentGenre(contentID string, genres []string) error {
	const op = "postgresql.SaveContentGenre"

	for _, genre := range genres {
		query := `INSERT INTO content_genres (content_id, genre) VALUES ($1, $2)`
		_, err := s.DB.Exec(query, contentID, genre)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
