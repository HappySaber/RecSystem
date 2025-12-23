package igdb

import (
	"catalog-microservice/internal/domain/models"
	"catalog-microservice/internal/storage/postgresql"
	"fmt"
)

type Importer struct {
	client  *Client
	content *postgresql.Storage
	games   *postgresql.Storage
}

func NewImporter(
	client *Client,
	content *postgresql.Storage,
	games *postgresql.Storage,
) *Importer {
	return &Importer{client, content, games}
}

// ImportPopularMovies импортирует первые 1000 популярных фильмов (по 20 на страницу)
func (i *Importer) ImportPopularMovies() error {
	const totalMovies = 1000
	const pageSize = 20
	totalPages := totalMovies / pageSize

	for page := 1; page <= totalPages; page++ {
		games, err := i.client.GetPopularMoviesPage(page)
		if err != nil {
			return fmt.Errorf("failed to fetch popular games page %d: %w", page, err)
		}

		for _, m := range games {
			m = m
			// raw, _ := json.Marshal(m)

			// content := &models.Content{
			// 	Type:           "movie",
			// 	ExternalSource: "tmdb",
			// 	ExternalID:     fmt.Sprint(m["id"]),
			// 	Title:          fmt.Sprint(m["title"]),
			// 	Description:    fmt.Sprint(m["overview"]),
			// 	PosterURL:      fmt.Sprint(m["poster_path"]),
			// 	ReleaseDate:    fmt.Sprint(m["release_date"]),
			// }

			// contentID, err := i.content.CreateOrUpdateContent(content)
			// if err != nil {
			// 	// Можно пропускать уже существующие записи
			// 	if err.Error() != "storage.posgresql.Create: unique constraint violation" {
			// 		return fmt.Errorf("failed to create content: %w", err)
			// 	}
			// }
			movie := &models.MovieDetails{}

			if err := i.games.SaveMovie(movie); err != nil {
				return fmt.Errorf("failed to save movie details: %w", err)
			}
		}
	}

	return nil
}
