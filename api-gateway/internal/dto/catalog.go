package dto

type Content struct {
	ID             string `json:"id"`
	ExternalSource string `json:"external_source"`
	ExternalID     string `json:"external_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	PosterURL      string `json:"poster_url"`
	ReleaseDate    string `json:"release_date"`
}
