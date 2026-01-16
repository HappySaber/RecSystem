package useractions

import (
	"log/slog"
	"time"
)

type UserAction struct {
	log       *slog.Logger
	uaTracker UserActionTracker
	tokenTTL  time.Duration
}

type UserActionTracker interface{}
