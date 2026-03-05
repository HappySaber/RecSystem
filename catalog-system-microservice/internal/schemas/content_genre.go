package schemas

type ContentGenre struct {
	ContentID string   `json:"content_id"`
	Genres    []string `json:"genres"`
}
