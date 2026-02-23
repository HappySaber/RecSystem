package dto

import "time"

type Content struct {
	ID             string `json:"id"`
	ExternalSource string `json:"external_source"`
	ExternalID     string `json:"external_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	PosterURL      string `json:"poster_url"`
	ReleaseDate    string `json:"release_date"`
}

type ContentShort struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type MovieDetails struct {
	ContentID     string `json:"content_id"`
	TmdbID        int32  `json:"tmdb_id"`
	OriginalTitle string `json:"original_title"`
	Runtime       *int32 `json:"runtime,omitempty"`
	Tagline       string `json:"tagline"`
	Status        string `json:"status"`
	Budget        int64  `json:"budget"`
	Revenue       int64  `json:"revenue"`
	Language      string `json:"language"`
}

type SeriesDetails struct {
	ContentID string

	TmdbID int32

	OriginalName     string
	Status           string
	FirstAirDate     *time.Time
	LastAirDate      *time.Time
	NumberOfSeasons  int32
	NumberOfEpisodes int32
	Language         string

	Genres      []byte
	Networks    []byte
	CastMembers []byte
	Images      []byte
	Videos      []byte

	RawData []byte
}

type GetContentRequest struct {
	ContentID string `json:"content_id"`
}
