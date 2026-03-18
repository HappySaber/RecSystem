package tmdb

import (
	"catalog-microservice/internal/domain/models"
	"catalog-microservice/internal/schemas"
	"catalog-microservice/internal/storage/postgresql"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// количество воркеров для параллельной загрузки деталей
	// TMDB бесплатный план — ~40 запросов/сек
	workerCount = 5
	// задержка между запросами одного воркера
	workerDelay = 150 * time.Millisecond
)

type Importer struct {
	client   *Client
	content  *postgresql.Storage
	movies   *postgresql.Storage
	series   *postgresql.Storage
	producer EventProducer
}

type EventProducer interface {
	Send(ctx context.Context, key string, value interface{}) error
}

func NewImporter(
	client *Client,
	content *postgresql.Storage,
	movies *postgresql.Storage,
	series *postgresql.Storage,
	producer EventProducer,
) *Importer {
	return &Importer{client, content, movies, series, producer}
}

func (i *Importer) ImportPopularMovies() error {
	const totalMovies = 1000
	const pageSize = 20
	totalPages := totalMovies / pageSize

	// собираем все ID фильмов сначала
	movieIDs := make([]int, 0, totalMovies)
	for page := 1; page <= totalPages; page++ {
		log.Printf("fetching page %d/%d", page, totalPages)

		moviesPage, err := i.client.GetPopularMoviesPage(page)
		if err != nil {
			return fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		for _, m := range moviesPage.Results {
			movieIDs = append(movieIDs, m.ID)
		}

		time.Sleep(250 * time.Millisecond)
	}

	log.Printf("collected %d movie IDs, starting import with %d workers", len(movieIDs), workerCount)

	// канал с ID фильмов для воркеров
	jobs := make(chan int, len(movieIDs))
	for _, id := range movieIDs {
		jobs <- id
	}
	close(jobs)

	var (
		wg        sync.WaitGroup
		processed atomic.Int64
		failed    atomic.Int64
	)

	// запускаем воркеры
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for movieID := range jobs {
				if err := i.importMovie(movieID); err != nil {
					log.Printf("worker %d: failed to import movie %d: %v", workerID, movieID, err)
					failed.Add(1)
				} else {
					count := processed.Add(1)
					if count%50 == 0 {
						log.Printf("progress: %d/%d movies imported", count, len(movieIDs))
					}
				}
				// задержка чтобы не превысить rate limit
				time.Sleep(workerDelay)
			}
		}(w)
	}

	wg.Wait()

	log.Printf("import completed: %d success, %d failed out of %d",
		processed.Load(), failed.Load(), len(movieIDs))

	return nil
}

// importMovie — импортирует один фильм
func (i *Importer) importMovie(movieID int) error {
	details, err := i.client.GetMovieDetails(movieID)
	if err != nil {
		return fmt.Errorf("fetch details: %w", err)
	}

	raw, _ := json.Marshal(details)

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
		return fmt.Errorf("save content: %w", err)
	}

	genres := extractStringsFromArray(details["genres"], "name")
	cast := []string{}
	crew := []string{}
	if credits, ok := details["credits"].(map[string]interface{}); ok {
		cast = extractStringsFromArray(credits["cast"], "name")
		crew = extractStringsFromArray(credits["crew"], "name")
	}

	movie := &models.MovieDetails{
		ContentID:     contentID,
		TmdbID:        int(getFloat64(details, "id")),
		OriginalTitle: getString(details, "original_title"),
		Runtime:       getIntPtr(details, "runtime"),
		Tagline:       getString(details, "tagline"),
		Status:        getString(details, "status"),
		Budget:        getInt64(details, "budget"),
		Revenue:       getInt64(details, "revenue"),
		Genres:        toJSON(genres),
		CastMembers:   toJSON(cast),
		Crew:          toJSON(crew),
		Images:        toJSON(details["images"]),
		Videos:        toJSON(details["videos"]),
		RawData:       raw,
	}

	if err := i.movies.SaveMovie(movie); err != nil {
		return fmt.Errorf("save movie: %w", err)
	}

	// отправляем событие в kafka
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	event := schemas.ContentGenre{
		ContentID: contentID,
		Genres:    genres,
	}

	if err := i.producer.Send(ctx, contentID, event); err != nil {
		// не фейлим импорт из-за kafka — просто логируем
		log.Printf("failed to publish kafka event for %s: %v", contentID, err)
	} else {
		log.Printf("kafka event published for %s", contentID)
	}

	return nil
}

func extractStringsFromArray(arr interface{}, key string) []string {
	result := []string{}
	items, ok := arr.([]interface{})
	if !ok || arr == nil {
		return result
	}
	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if val, exists := obj[key]; exists && val != nil {
			if str, ok := val.(string); ok {
				result = append(result, str)
			}
		}
	}
	return result
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprint(v)
	}
	return ""
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

func getInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			return int64(f)
		}
	}
	return 0
}

func getIntPtr(m map[string]interface{}, key string) *int {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			val := int(f)
			return &val
		}
	}
	return nil
}

func toJSON(v interface{}) json.RawMessage {
	if v == nil {
		return json.RawMessage(`[]`)
	}
	b, _ := json.Marshal(v)
	return b
}
