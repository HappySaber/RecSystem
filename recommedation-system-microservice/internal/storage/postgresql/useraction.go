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

	query := `
	INSERT INTO user_content_action
	    (user_id, content_id, action_id, rating, duration_sec)
	SELECT
	    $1,
	    $2,
	    a.id,
	    $4,
	    $5
	FROM actions a
	WHERE a.code = $3
	`

	_, err := s.DB.ExecContext(
		ctx,
		query,
		event.UserID,
		event.ContentID,
		event.Action,
		event.Rating,
		event.DurationSec,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserLikedContent(
	ctx context.Context,
	userID string) ([]string, error) {
	const op = "postgresql.GetUserLikedContent"

	query := `SELECT content_id FROM user_actions WHERE user_id = $1 AND action_type = 'like'`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var contentIDs []string
	for rows.Next() {
		var contentID string
		if err := rows.Scan(&contentID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		contentIDs = append(contentIDs, contentID)

	}
	return contentIDs, nil
}

// func (s *Storage) SaveTrackActions(
// 	ctx context.Context,
// 	events []models.UserActionEvent,
// ) error {
// 	//TODO realisation of batch insert
// 	return nil
// }
