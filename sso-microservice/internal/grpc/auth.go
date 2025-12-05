package auth

import (
	"context"
	"errors"
	sso1 "sso-microservice/internal/pb/sso.v1"
	"sso-microservice/internal/services/auth"
	"sso-microservice/internal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int,
	) (accessToken string, err error)
	//refreshToken string,
	RegisterNewUser(ctx context.Context,
		email string,
		name string,
		surname string,
		role string,
		password string,
	) (userID string, err error)
	UserRole(ctx context.Context,
		userID string,
	) (isAdmin string, err error)
	Logout(ctx context.Context,
		refToken string,
	) (isLogouted bool, err error)
}

type serverAPI struct {
	sso1.UnimplementedAuthServiceServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	sso1.RegisterAuthServiceServer(gRPC, &serverAPI{auth: auth})

}

func (s *serverAPI) Login(ctx context.Context, req *sso1.LoginRequest) (*sso1.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	accessToken, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso1.LoginResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso1.RegisterRequest) (*sso1.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetName(), req.GetSurname(), req.GetRole(), req.GetPassword())

	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso1.RegisterResponse{
		UserId: userID,
	}, nil
}
func (s *serverAPI) Logout(ctx context.Context, req *sso1.LogoutRequest) (*sso1.LogoutResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong argument")
	}

	isLogouted, err := s.auth.Logout(ctx, req.GetRefreshToken())
	if err != nil {

		return nil, status.Error(codes.Internal, "Coudnot logout")
	}

	return &sso1.LogoutResponse{
		Success: isLogouted,
		Message: "User logged out",
	}, nil
}
func (s *serverAPI) IsAdmin(ctx context.Context, req *sso1.IsAdminRequest) (*sso1.IsAdminResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	role, err := s.auth.UserRole(ctx, req.GetUserId())

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	if role == "admin" {
		return &sso1.IsAdminResponse{
			IsAdmin: true,
		}, nil
	}

	return &sso1.IsAdminResponse{
		IsAdmin: false,
	}, nil
}
func (s *serverAPI) RefreshToken(ctx context.Context, req *sso1.RefreshTokenRequest) (*sso1.RefreshTokenResponse, error) {
	panic("implement me")
}
