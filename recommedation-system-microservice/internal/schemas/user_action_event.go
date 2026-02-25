package schemas

type UserActionEvent struct {
	UserID    string `json:"user_id"`
	ContentID string `json:"content_id"`
	Action    string `json:"action"`
	Rating    *int32 `json:"rating,omitempty"`
	Duration  *int32 `json:"duration_sec,omitempty"`
}
