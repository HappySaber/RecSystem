package tmdb

import (
	"catalog-microservice/internal/domain/models"
	"encoding/json"
	"fmt"
)

func (i *Importer) ImportPopularTVShows() error {
	const pages = 50 // 1000 сериалов

	for page := 1; page <= pages; page++ {
		shows, err := i.client.GetPopularTVPage(page)
		if err != nil {
			return err
		}

		for _, s := range shows {
			raw, _ := json.Marshal(s)

			content := &models.Content{
				Type:           "series",
				ExternalSource: "tmdb",
				ExternalID:     fmt.Sprint(s.ID),
				Title:          s.Name,
				Description:    s.Overview,
				PosterURL:      s.PosterPath,
				ReleaseDate:    s.FirstAirDate,
			}

			contentID, err := i.content.CreateOrUpdateContent(content)
			if err != nil {
				return err
			}

			if err := i.series.SaveSeries(contentID, s.ID, string(raw)); err != nil {
				return err
			}
		}
	}
	return nil
}
