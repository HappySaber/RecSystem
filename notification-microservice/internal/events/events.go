package events

import "time"

type UserRegisteredEvent struct {
	EventID   string    `json:"eventId"`
	UserID    string    `json:"userId"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	CreatedAt time.Time `json:"createdAt"`
}
