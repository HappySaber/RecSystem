package recommendation

import (
	"context"

	recs "rec-system-microservice/internal/pb/recommendation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecommendationEngine interface {
	GetRecommendations(
		ctx context.Context,
		userID string,
		limit int,
	) ([]string, error)

	GetRecommendationsByGenres(
		ctx context.Context,
		genres []string,
		limit int,
	) ([]string, error)

	GetSimilarContent(
		ctx context.Context,
		contentID string,
		limit int,
	) ([]string, error)

	GetTrendingContent(
		ctx context.Context,
		limit int,
	) ([]string, error)

	GetPopularContent(
		ctx context.Context,
		limit int,
	) ([]string, error)
}

//
// core
//

func (s serverAPI) GetRecommendations(
	ctx context.Context,
	req *recs.GetRecommendationsRequest,
) (*recs.GetRecommendationsResponse, error) {

	if req.GetUserId() == "" || req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}
	recIDs, err := s.engine.GetRecommendations(ctx, req.GetUserId(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.GetRecommendationsResponse{
		ContentIds: recIDs,
	}, nil
}

func (s serverAPI) GetRecommendationsByGenres(
	ctx context.Context,
	req *recs.GetRecommendationsByGenresRequest,
) (*recs.GetRecommendationsByGenresResponse, error) {

	if len(req.GetGenres()) == 0 || req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	recIDs, err := s.engine.GetRecommendationsByGenres(ctx, req.GetGenres(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.GetRecommendationsByGenresResponse{
		ContentIds: recIDs,
	}, nil
}

func (s serverAPI) GetSimilarContent(
	ctx context.Context,
	req *recs.GetSimilarContentRequest,
) (*recs.GetSimilarContentResponse, error) {

	if req.GetContentId() == "" || req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	recIDs, err := s.engine.GetSimilarContent(ctx, req.GetContentId(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.GetSimilarContentResponse{
		ContentIds: recIDs,
	}, nil
}

func (s serverAPI) GetTrendingContent(
	ctx context.Context,
	req *recs.GetTrendingContentRequest,
) (*recs.GetTrendingContentResponse, error) {

	if req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	recIDs, err := s.engine.GetTrendingContent(ctx, int(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.GetTrendingContentResponse{
		ContentIds: recIDs,
	}, nil
}

func (s serverAPI) GetPopularContent(
	ctx context.Context,
	req *recs.GetPopularContentRequest,
) (*recs.GetPopularContentResponse, error) {

	if req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	recIDs, err := s.engine.GetPopularContent(ctx, int(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.GetPopularContentResponse{
		ContentIds: recIDs,
	}, nil
}
