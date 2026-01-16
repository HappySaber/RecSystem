package recommendation

import (
	"log/slog"
	"time"
)

type Recommendation struct {
	log      *slog.Logger
	engine   RecommendationProvider
	tokenTTL time.Duration
}

type RecommendationProvider interface{}
