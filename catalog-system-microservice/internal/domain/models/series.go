package models

import "time"

type SeriesDetails struct {
	ContentID string

	TmdbID int

	// Основные метаданные сериала
	OriginalName     string
	Status           string // Ended, Returning Series, etc.
	FirstAirDate     *time.Time
	LastAirDate      *time.Time
	NumberOfSeasons  int
	NumberOfEpisodes int
	Language         string

	// JSONB-поля
	Genres      []byte // ["Drama", "Crime"]
	Networks    []byte // Netflix, HBO и т.д.
	CastMembers []byte
	Images      []byte
	Videos      []byte

	RawData []byte
}
