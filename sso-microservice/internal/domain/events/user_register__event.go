package events

import "time"

type UserRegisteredEvent struct {
	EventID   string
	UserID    string
	Email     string
	Name      string
	Surname   string
	CreatedAt time.Time
}
