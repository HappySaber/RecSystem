package models

// type MovieDetails struct {
// 	ContentID string
// 	TmdbID    int
// 	RawJSON   string
// }

type MovieDetails struct {
	ContentID string

	TmdbID        int
	OriginalTitle string
	Runtime       *int
	Tagline       string
	Status        string
	Budget        int64
	Revenue       int64
	Language      string

	Genres      []byte
	CastMembers []byte
	Crew        []byte
	Images      []byte
	Videos      []byte
	RawData     []byte
}
