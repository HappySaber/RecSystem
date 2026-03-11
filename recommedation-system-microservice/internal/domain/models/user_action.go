package models

import "time"

type ActionType string

const (
	ActionView     ActionType = "VIEW"
	ActionLike     ActionType = "LIKE"
	ActionDislike  ActionType = "DISLIKE"
	ActionRate     ActionType = "RATE"
	ActionFavorite ActionType = "FAVORITE"
)

type UserActionEvent struct {
	UserID    string     `json:"user_id"`
	ContentID string     `json:"content_id"`
	Action    ActionType `json:"action"`
	Timestamp time.Time  `json:"timestamp"`

	Rating      *int32 `json:"rating,omitempty"`
	DurationSec *int32 `json:"duration_sec,omitempty"`
}
