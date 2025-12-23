package tmdb

import (
	"catalog-microservice/internal/domain/models"
	"catalog-microservice/internal/storage/postgresql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Importer struct {
	client  *Client
	content *postgresql.Storage
	movies  *postgresql.Storage
	series  *postgresql.Storage
}

func NewImporter(
	client *Client,
	content *postgresql.Storage,
	movies *postgresql.Storage,
	series *postgresql.Storage,
) *Importer {
	return &Importer{client, content, movies, series}
}

func extractStringsFromArray(arr interface{}, key string) []string {
	result := []string{}
	items, ok := arr.([]interface{})
	if !ok || arr == nil {
		return result
	}
	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok || obj == nil {
			continue
		}
		if val, exists := obj[key]; exists && val != nil {
			if str, ok := val.(string); ok {
				result = append(result, str) // напрямую строка, без fmt.Sprint
			}
		}
	}
	return result
}

// ImportPopularMovies imports the first 1000 popular movies from TMDb
func (i *Importer) ImportPopularMovies() error {
	const totalMovies = 1000
	const pageSize = 20
	totalPages := totalMovies / pageSize

	for page := 1; page <= totalPages; page++ {
		moviesPage, err := i.client.GetPopularMoviesPage(page)
		if err != nil {
			return fmt.Errorf("failed to fetch popular movies page %d: %w", page, err)
		}

		for _, m := range moviesPage.Results {
			// Get full movie details
			details, err := i.client.GetMovieDetails(m.ID)
			if err != nil {
				return fmt.Errorf("failed to fetch movie details %d: %w", m.ID, err)
			}

			raw, _ := json.Marshal(details)

			// Create content
			content := &models.Content{
				Type:           "movie",
				ExternalSource: "tmdb",
				ExternalID:     fmt.Sprint(details["id"]),
				Title:          fmt.Sprint(details["title"]),
				Description:    fmt.Sprint(details["overview"]),
				PosterURL:      fmt.Sprint(details["poster_path"]),
				ReleaseDate:    fmt.Sprint(details["release_date"]),
			}

			contentID, err := i.content.CreateOrUpdateContent(content)
			if err != nil {
				return fmt.Errorf("failed to create/update content: %w", err)
			}

			// Extract genres as array of maps [{id: ..., name: ...}]
			genres := extractStringsFromArray(details["genres"], "name")
			cast := []string{}
			crew := []string{}
			if credits, ok := details["credits"].(map[string]interface{}); ok {
				cast = extractStringsFromArray(credits["cast"], "name")
				crew = extractStringsFromArray(credits["crew"], "name")
			}

			// Save movie details
			movie := &models.MovieDetails{
				ContentID:     contentID,
				TmdbID:        m.ID,
				OriginalTitle: getString(details, "original_title"),
				Runtime:       getIntPtr(details, "runtime"),
				Tagline:       getString(details, "tagline"),
				Status:        getString(details, "status"),
				Budget:        getInt64(details, "budget"),
				Revenue:       getInt64(details, "revenue"),
				Language:      getString(details, "original_language"),

				Genres:      toJSON(genres),
				CastMembers: toJSON(cast),
				Crew:        toJSON(crew),
				Images:      toJSON(details["images"]),
				Videos:      toJSON(details["videos"]),
				RawData:     raw,
			}

			log.Println(string(movie.Genres))

			if err := i.movies.SaveMovie(movie); err != nil {
				return fmt.Errorf("failed to save movie details: %w", err)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprint(v)
	}
	return ""
}

func getInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			return int64(f)
		}
	}
	return 0
}

func toJSON(v interface{}) json.RawMessage {
	if v == nil {
		return json.RawMessage(`[]`)
	}
	b, _ := json.Marshal(v)
	return b
}

func getIntPtr(m map[string]interface{}, key string) *int {
	if v, ok := m[key]; ok && v != nil {
		if i, ok := v.(float64); ok {
			val := int(i)
			return &val
		}
	}
	return nil
}
