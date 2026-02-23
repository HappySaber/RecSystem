package dto

type TrackUserActionRequest struct {
	UserID    string `json:"user_id"`
	ContentID string `json:"content_id"`
	Action    string `json:"action"`
	Rating    *int32 `json:"rating"`
	Duration  *int32 `json:"duration_sec"`
}

type TrackUserActionResponse struct{}
