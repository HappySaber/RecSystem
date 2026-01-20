package recommendation

import recs "rec-system-microservice/internal/pb/recommendation"

type serverAPI struct {
	recs.UnimplementedRecommendationServiceServer
	engine   RecommendationEngine
	prefs    UserPreferencesService
	actions  UserActionTracker
	aiEngine AIRecommendationEngine
}
