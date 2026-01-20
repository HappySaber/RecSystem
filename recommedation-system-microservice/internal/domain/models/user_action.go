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

type ActionMeta struct {
	Rating      *int32
	DurationSec *int32
	Timestamp   time.Time
}

type UserActionEvent struct {
	UserID    string
	ContentID string
	Action    ActionType
	Meta      ActionMeta
}
