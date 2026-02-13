package client

import (
	pbSSO "api-gateway/internal/pbs/sso/sso"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type SSOClient struct {
	client pbSSO.AuthServiceClient
}

func NewSSOClient(addr string) (*SSOClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sso service: %w", err)
	}

	c := pbSSO.NewAuthServiceClient(conn)
	return &SSOClient{
		client: c,
	}, nil
}

func (c *SSOClient) Register(ctx context.Context, email, password, name, surname, role string) (string, error) {
	resp, err := c.client.Register(
		ctx,
		&pbSSO.RegisterRequest{
			Email:    email,
			Password: password,
			Name:     name,
			Surname:  surname,
			Role:     role,
		},
	)
	if err != nil {
		return "", err
	}

	return resp.UserId, nil
}

func (c *SSOClient) Login(ctx context.Context, email, password string) (string, error) {
	resp, err := c.client.Login(
		ctx,
		&pbSSO.LoginRequest{
			Email:    email,
			Password: password,
		},
	)
	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}
