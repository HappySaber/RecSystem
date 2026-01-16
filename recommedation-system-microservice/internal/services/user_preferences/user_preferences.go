package userpreferences

import (
	"context"
	"log/slog"
	"time"
)

type UserPreferences struct {
	log             *slog.Logger
	userPreferences UserPreferencesService
	tokenTTL        time.Duration
}

type UserPreferencesService interface {
	GetUserPreferences(ctx context.Context, userID string) ([]string, error)
	SetUserPreferences(ctx context.Context, userID string, genres []string) error
	ResetUserPreferences(ctx context.Context, userID string) error
}
