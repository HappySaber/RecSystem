package recommendation

import (
	"context"
	recs "rec-system-microservice/internal/pb/recommendation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AIRecommendationEngine interface {
	GetImplicitRecommendations(
		ctx context.Context,
		userID string,
		limit int,
	) ([]string, error)

	GetExplicitRecommendations(
		ctx context.Context,
		query string,
		limit int,
	) ([]string, error)
}

func (s *serverAPI) GetImplicitAIRecommendations(
	ctx context.Context,
	req *recs.GetAIRecommendationsRequest,
) (*recs.GetAIRecommendationsResponse, error) {

	if req.GetUserId() == "" || req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}
	implicitRecommendations, err := s.aiEngine.GetImplicitRecommendations(ctx, req.GetUserId(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &recs.GetAIRecommendationsResponse{
		ContentIds: implicitRecommendations,
	}, nil
}

func (s *serverAPI) GetExplicitAIRecommendations(
	ctx context.Context,
	req *recs.GetAIRecommendationsRequest,
) (*recs.GetAIRecommendationsResponse, error) {

	if req.GetQuery() == "" || req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}
	explicitRecommendations, err := s.aiEngine.GetExplicitRecommendations(ctx, req.GetQuery(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &recs.GetAIRecommendationsResponse{
		ContentIds: explicitRecommendations,
	}, nil
}
