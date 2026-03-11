package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/lib/pq"
)

func (s *Storage) GetUserRecommendations(ctx context.Context, userID string, limit int) ([]string, error) {
	const op = "storage.GetUserRecommendations"

	query := `
		SELECT
    cg.content_id,
    SUM(pref.score) AS total_score
FROM content_genres cg
JOIN genres g ON cg.genre_id = g.id
JOIN (
    SELECT
        g.name,
        SUM(a.weight) AS score
    FROM user_content_action uca
    JOIN actions a ON uca.action_id = a.id
    JOIN content_genres cg ON uca.content_id = cg.content_id
    JOIN genres g ON cg.genre_id = g.id
    WHERE uca.user_id = $1
    GROUP BY g.name
) pref ON pref.name = g.name
WHERE cg.content_id NOT IN (
    SELECT content_id
    FROM user_content_action
    WHERE user_id = $1
)
GROUP BY cg.content_id
ORDER BY total_score DESC
LIMIT $2;
	`

	rows, err := s.DB.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var contentID string
		var score float64

		if err := rows.Scan(&contentID, &score); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		result = append(result, contentID)
	}
	log.Println(result)

	return result, nil
}

func (s *Storage) GetRecommendationsByGenres(ctx context.Context, genres []string, limit int) ([]string, error) {
	const op = "storage.GetRecommendationsByGenres"

	query := `
	SELECT cg.content_id
	FROM content_genres cg
	JOIN genres g ON cg.genre_id = g.id
	WHERE g.name = ANY($1)
	GROUP BY cg.content_id
	ORDER BY COUNT(*) DESC
	LIMIT $2
	`

	rows, err := s.DB.QueryContext(ctx, query, pq.Array(genres), limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var contentID string
		if err := rows.Scan(&contentID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, contentID)
	}

	return result, nil
}

func (s *Storage) GetSimilarContent(ctx context.Context, contentID string, limit int) ([]string, error) {
	const op = "storage.GetSimilarContent"

	query := `
	SELECT cg2.content_id
	FROM content_genres cg1
	JOIN content_genres cg2 ON cg1.genre_id = cg2.genre_id
	WHERE cg1.content_id = $1
	  AND cg2.content_id != $1
	GROUP BY cg2.content_id
	ORDER BY COUNT(*) DESC
	LIMIT $2
	`

	rows, err := s.DB.QueryContext(ctx, query, contentID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var cid string
		if err := rows.Scan(&cid); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, cid)
	}

	return result, nil
}

func (s *Storage) GetTrendingContent(ctx context.Context, limit int) ([]string, error) {
	const op = "storage.GetTrendingContent"

	query := `
	SELECT uca.content_id
	FROM user_content_action uca
	JOIN actions a ON uca.action_id = a.id
	WHERE a.code = 'VIEW' AND uca.created_at > NOW() - INTERVAL '24 hours'
	GROUP BY uca.content_id
	ORDER BY COUNT(*) DESC
	LIMIT $1
	`

	rows, err := s.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var contentID string
		if err := rows.Scan(&contentID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, contentID)
	}

	return result, nil
}

func (s *Storage) GetPopularContent(ctx context.Context, limit int) ([]string, error) {
	const op = "storage.GetPopularContent"

	query := `
	SELECT content_id
	FROM content_popularity_agg
	ORDER BY score DESC
	LIMIT $1
	`

	rows, err := s.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var contentID string
		if err := rows.Scan(&contentID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, contentID)
	}

	return result, nil
}
