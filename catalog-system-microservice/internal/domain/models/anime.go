package models

type AnimeDetails struct {
	ContentID string

	AniListID int
	MALID     *int

	OriginalTitle string
	Format        string
	Status        string
	Season        string
	SeasonYear    *int

	EpisodesCount   *int
	EpisodeDuration *int

	StartDate *string
	EndDate   *string

	Language string

	Genres      []byte
	Tags        []byte
	Studios     []byte
	Characters  []byte
	VoiceActors []byte

	MeanScore  *int
	Popularity int
	Favourites int

	TrailerURL *string

	RawData []byte
}
