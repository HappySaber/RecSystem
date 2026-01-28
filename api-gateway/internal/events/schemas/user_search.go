package schemas

import "time"

type UserSearchEvent struct {
	EventID   string    `json:"event_id"`
	RequestID string    `json:"request_id"`
	UserID    string    `json:"user_id"`
	Query     string    `json:"query"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}
