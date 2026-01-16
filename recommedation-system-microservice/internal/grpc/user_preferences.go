package recommendation

import (
	"context"
	recs "rec-system-microservice/internal/pb/recommendation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserPreferencesService interface {
	// preferences
	GetUserPreferences(
		ctx context.Context,
		userID string,
	) ([]string, error)

	SetUserPreferences(
		ctx context.Context,
		userID string,
		genres []string,
	) error

	ResetUserPreferences(
		ctx context.Context,
		userID string,
	) error

	RebuildUserPreferences(
		ctx context.Context,
		userID string,
	) error
}

func (s serverAPI) GetUserPreferences(
	ctx context.Context,
	req *recs.GetUserPreferencesRequest,
) (*recs.GetUserPreferencesResponse, error) {

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	genres, err := s.prefs.GetUserPreferences(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.GetUserPreferencesResponse{
		PreferredGenres: genres,
	}, nil
}

func (s serverAPI) SetUserPreferences(
	ctx context.Context,
	req *recs.SetUserPreferencesRequest,
) (*recs.SetUserPreferencesResponse, error) {

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	err := s.prefs.SetUserPreferences(ctx, req.GetUserId(), req.GetPreferredGenres())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.SetUserPreferencesResponse{
		UserPreferencesSetted: true,
	}, nil
}

func (s serverAPI) ResetUserPreferences(
	ctx context.Context,
	req *recs.ResetUserPreferencesRequest,
) (*recs.ResetUserPreferencesResponse, error) {

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	err := s.prefs.ResetUserPreferences(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.ResetUserPreferencesResponse{
		UserPreferencesResetted: true,
	}, nil
}

func (s serverAPI) RebuildUserPreferences(
	ctx context.Context,
	req *recs.RebuildUserPreferencesRequest,
) (*recs.RebuildUserPreferencesResponce, error) {

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	err := s.prefs.RebuildUserPreferences(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.RebuildUserPreferencesResponce{}, nil
}
