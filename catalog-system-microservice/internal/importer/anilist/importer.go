package anilist

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
	anime   *postgresql.Storage
}

func NewImporter(
	client *Client,
	content *postgresql.Storage,
	anime *postgresql.Storage,
) *Importer {
	return &Importer{client, content, anime}
}

func (i *Importer) ImportPopularAnime() error {
	const pages = 50 // 50 × 20 = 1000

	for page := 1; page <= pages; page++ {
		list, err := i.client.GetPopularAnimePage(page)
		if err != nil {
			return err
		}

		for _, a := range list {
			raw, _ := json.Marshal(a)

			title := a.Title.English
			if title == "" {
				title = a.Title.Romaji
			}

			content := &models.Content{
				Type:           "anime",
				ExternalSource: "anilist",
				ExternalID:     fmt.Sprint(a.ID),
				Title:          title,
				Description:    a.Description,
				PosterURL:      a.CoverImage.Large,
			}

			contentID, err := i.content.CreateOrUpdateContent(content)
			if err != nil {
				return err
			}

			// ---- списки ----
			tags := []string{}
			for _, t := range a.Tags {
				tags = append(tags, t.Name)
			}

			studios := []string{}
			for _, s := range a.Studios.Nodes {
				studios = append(studios, s.Name)
			}

			characters := []string{}
			for _, c := range a.Characters.Nodes {
				characters = append(characters, c.Name.Full)
			}

			voiceActors := []string{}
			for _, v := range a.VoiceActors.Nodes {
				voiceActors = append(voiceActors, v.Name.Full)
			}

			// ---- даты ----
			startDate := dateToString(a.StartDate)
			endDate := dateToString(a.EndDate)

			// ---- трейлер ----
			var trailerURL *string
			if a.Trailer != nil && a.Trailer.Site == "youtube" {
				url := "https://youtube.com/watch?v=" + a.Trailer.ID
				trailerURL = &url
			}

			anime := &models.AnimeDetails{
				ContentID: contentID,

				AniListID: a.ID,
				MALID:     a.IDMal,

				OriginalTitle: title,
				Format:        a.Format,
				Status:        a.Status,
				Season:        a.Season,
				SeasonYear:    a.SeasonYear,

				EpisodesCount:   a.Episodes,
				EpisodeDuration: a.Duration,

				StartDate: startDate,
				EndDate:   endDate,

				Language: "ja",

				Genres:      toJSON(a.Genres),
				Tags:        toJSON(tags),
				Studios:     toJSON(studios),
				Characters:  toJSON(characters),
				VoiceActors: toJSON(voiceActors),

				MeanScore:  a.MeanScore,
				Popularity: a.Popularity,
				Favourites: a.Favourites,

				TrailerURL: trailerURL,
				RawData:    raw,
			}

			log.Println("anime:", anime.OriginalTitle)

			if err := i.anime.SaveAnime(anime); err != nil {
				return err
			}

			time.Sleep(150 * time.Millisecond)
		}
	}

	return nil
}

func toJSON(v interface{}) []byte {
	if v == nil {
		return []byte(`[]`)
	}
	b, _ := json.Marshal(v)
	return b
}

func dateToString(d *DateDTO) *string {
	if d == nil || d.Year == nil || d.Month == nil || d.Day == nil {
		return nil
	}
	s := fmt.Sprintf("%04d-%02d-%02d", *d.Year, *d.Month, *d.Day)
	return &s
}
