package tmdb

import (
	"catalog-microservice/internal/domain/models"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (i *Importer) ImportPopularSeries() error {
	const totalSeries = 1000
	const pageSize = 20
	totalPages := totalSeries / pageSize

	for page := 1; page <= totalPages; page++ {
		seriesPage, err := i.client.GetPopularSeriesPage(page)
		if err != nil {
			return fmt.Errorf("failed to fetch popular series page %d: %w", page, err)
		}

		for _, s := range seriesPage.Results {
			details, err := i.client.GetSeriesDetails(s.ID)
			if err != nil {
				return fmt.Errorf("failed to fetch series details %d: %w", s.ID, err)
			}

			raw, _ := json.Marshal(details)

			// content (общая таблица)
			content := &models.Content{
				Type:           "series",
				ExternalSource: "tmdb",
				ExternalID:     fmt.Sprint(details["id"]),
				Title:          getString(details, "name"),
				Description:    getString(details, "overview"),
				PosterURL:      getString(details, "poster_path"),
				ReleaseDate:    getString(details, "first_air_date"),
			}

			contentID, err := i.content.CreateOrUpdateContent(content)
			if err != nil {
				return fmt.Errorf("failed to create/update content: %w", err)
			}

			// genres
			genres := extractStringsFromArray(details["genres"], "name")

			// cast
			cast := []string{}
			if credits, ok := details["credits"].(map[string]interface{}); ok {
				cast = extractStringsFromArray(credits["cast"], "name")
			}

			series := &models.SeriesDetails{
				ContentID: contentID,
				TmdbID:    s.ID,

				OriginalName: getString(details, "original_name"),
				Status:       getString(details, "status"),

				FirstAirDate: parseDate(details["first_air_date"]),
				LastAirDate:  parseDate(details["last_air_date"]),

				NumberOfSeasons:  getInt(details, "number_of_seasons"),
				NumberOfEpisodes: getInt(details, "number_of_episodes"),
				Language:         getString(details, "original_language"),

				Genres:      toJSON(genres),
				Networks:    toJSON(details["networks"]),
				CastMembers: toJSON(cast),
				Images:      toJSON(details["images"]),
				Videos:      toJSON(details["videos"]),

				RawData: raw,
			}

			log.Println("series:", series.OriginalName)
			if err := i.series.SaveSeries(series); err != nil {
				return fmt.Errorf("failed to save series details: %w", err)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return 0
}

func parseDate(value interface{}) *time.Time {
	str, ok := value.(string)
	if !ok || str == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return nil
	}

	return &t
}
