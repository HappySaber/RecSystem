package dto

type GetExplicitRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

type GetExplicitResponse struct {
	Titles     []string `json:"titles"`
	ContentIDs []string `json:"content_ids,omitempty"` // deprecated, для совместимости
}

type GetRecommendationsRequest struct {
	Limit int `json:"limit"`
}

type GetRecommendationsResponse struct {
	ContentIDs []string `json:"content_ids"`
}

type GetRecommendationsByGenresRequest struct {
	Genres []string `json:"genres"`
	Limit  int      `json:"limit"`
}

type GetRecommendationsByGenresResponse struct {
	ContentIDs []string `json:"content_ids"`
}

type GetSimilarContentRequest struct {
	ContentID string `json:"content_id"`
	Limit     int    `json:"limit"`
}

type GetSimilarContentResponse struct {
	ContentIDs []string `json:"content_ids"`
}

type GetTrendingContentRequest struct {
	Limit int `json:"limit"`
}

type GetTrendingContentResponse struct {
	ContentIDs []string `json:"content_ids"`
}

type GetPopularContentRequest struct {
	Limit int `json:"limit"`
}

type GetPopularContentResponse struct {
	ContentIDs []string `json:"content_ids"`
}
