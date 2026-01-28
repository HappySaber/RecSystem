package dto

type GetExplicitRequest struct {
	UserID string `json:"user_id"`
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
}

type GetExplicitResponse struct {
	ContentIDs []string `json:"content_ids"`
}
