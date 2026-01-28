package schemas

import "time"

type UserActionEvent struct {
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	ContentID string    `json:"content_id"`
	Action    string    `json:"action"`
	Rating    *int32    `json:"rating,omitempty"`
	Duration  *int32    `json:"duration_sec,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}
