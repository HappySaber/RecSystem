package auth

import (
	"context"
	sso1 "sso-microservice/internal/pb/sso.v1"

	"google.golang.org/grpc"
)

type serverAPI struct {
	sso1.UnimplementedAuthServiceServer
}

func Register(gRPC *grpc.Server) {
	sso1.RegisterAuthServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(cxt context.Context, req *sso1.LoginRequest) (*sso1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(cxt context.Context, req *sso1.RegisterRequest) (*sso1.RegisterResponse, error) {
	panic("implement me")
}
func (s *serverAPI) Logout(cxt context.Context, req *sso1.LogoutRequest) (*sso1.LogoutResponse, error) {
	panic("implement me")
}
func (s *serverAPI) IsAdmin(cxt context.Context, req *sso1.IsAdminRequest) (*sso1.IsAdminResponse, error) {
	panic("implement me")
}
func (s *serverAPI) RefreshToken(cxt context.Context, req *sso1.RefreshTokenRequest) (*sso1.RefreshTokenResponse, error) {
	panic("implement me")
}
