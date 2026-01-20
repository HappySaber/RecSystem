package recommendation

import (
	recs "rec-system-microservice/internal/pb/recommendation"

	"google.golang.org/grpc"
)

type serverAPI struct {
	recs.UnimplementedRecommendationServiceServer
	engine   RecommendationEngine
	prefs    UserPreferencesService
	actions  UserActionTracker
	aiEngine AIRecommendationEngine
}

func Register(gRPC *grpc.Server, engine RecommendationEngine, prefs UserPreferencesService, actions UserActionTracker, aiEngine AIRecommendationEngine) {
	recs.RegisterRecommendationServiceServer(gRPC, &serverAPI{
		engine:   engine,
		prefs:    prefs,
		actions:  actions,
		aiEngine: aiEngine,
	})

}
