package postgresql

import (
	"context"
	"fmt"
	"rec-system-microservice/internal/domain/models"
)

func (s *Storage) SaveTrackAction(
	ctx context.Context,
	event models.UserActionEvent,
) error {
	const op = "postgresql.SaveTrackAction"

	query := `INSERT INTO user_actions
	(user_id, content_id, action_type, rating, duration_sec, created_at)
	VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err := s.db.ExecContext(
		ctx,
		query,
		event.UserID,
		event.ContentID,
		event.Action,
		event.Meta.Rating,
		event.Meta.DurationSec,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
