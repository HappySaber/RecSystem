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

type UserPreferencesService interface{}

func (ur *UserPreferences) GetUserPreferences(
	ctx context.Context,
	userID string,
) ([]string, error) {
	return []string{}, nil
}

func (ur *UserPreferences) SetUserPreferences(
	ctx context.Context,
	userID string,
	genres []string,
) error {
	return nil
}
func (ur *UserPreferences) ResetUserPreferences(
	ctx context.Context,
	userID string,
) error {
	return nil
}

func (ur *UserPreferences) RebuildUserPreferences(
	ctx context.Context,
	userID string,
) error {
	return nil
}
