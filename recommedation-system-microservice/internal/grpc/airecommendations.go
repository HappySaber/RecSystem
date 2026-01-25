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

// func (s *serverAPI) GetImplicitAIRecommendations(
// 	ctx context.Context,
// 	req *recs.GetAIRecommendationsRequest,
// ) (*recs.GetAIRecommendationsResponse, error) {

// 	if req.GetUserId() == "" || req.GetLimit() <= 0 {
// 		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
// 	}
// 	implicitRecommendations, err := s.aiEngine.GetImplicitRecommendations(ctx, req.GetUserId(), int(req.GetLimit()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &recs.GetAIRecommendationsResponse{
// 		ContentIds: implicitRecommendations,
// 	}, nil
// }

// func (s *serverAPI) GetExplicitAIRecommendations(
// 	ctx context.Context,
// 	req *recs.GetAIRecommendationsRequest,
// ) (*recs.GetAIRecommendationsResponse, error) {

// 	if req.GetQuery() == "" || req.GetLimit() <= 0 {
// 		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
// 	}
// 	explicitRecommendations, err := s.aiEngine.GetExplicitRecommendations(ctx, req.GetQuery(), int(req.GetLimit()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &recs.GetAIRecommendationsResponse{
// 		ContentIds: explicitRecommendations,
// 	}, nil
// }

func (s *serverAPI) GetAIRecommendations(
	ctx context.Context,
	req *recs.GetAIRecommendationsRequest,
) (*recs.GetAIRecommendationsResponse, error) {

	if req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be > 0")
	}

	switch req.GetMode() {

	case recs.RecommendationMode_IMPLICIT:
		if req.GetUserId() == "" {
			return nil, status.Error(codes.InvalidArgument, "user_id is required for IMPLICIT mode")
		}

		recsIDs, err := s.aiEngine.GetImplicitRecommendations(
			ctx,
			req.GetUserId(),
			int(req.GetLimit()),
		)
		if err != nil {
			return nil, err
		}

		return &recs.GetAIRecommendationsResponse{
			ContentIds: recsIDs,
		}, nil

	case recs.RecommendationMode_EXPLICIT:
		if req.GetQuery() == "" {
			return nil, status.Error(codes.InvalidArgument, "query is required for EXPLICIT mode")
		}

		recsIDs, err := s.aiEngine.GetExplicitRecommendations(
			ctx,
			req.GetQuery(),
			int(req.GetLimit()),
		)
		if err != nil {
			return nil, err
		}

		return &recs.GetAIRecommendationsResponse{
			ContentIds: recsIDs,
		}, nil

	default:
		return nil, status.Error(codes.InvalidArgument, "unknown recommendation mode")
	}
}
