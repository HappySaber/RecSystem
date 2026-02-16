package client

import (
	"api-gateway/internal/dto"
	pbCatalog "api-gateway/internal/pbs/catalog"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type CatalogClient struct {
	client pbCatalog.CatalogServiceClient
}

func NewCatalogClient(addr string) (*CatalogClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sso service: %w", err)
	}

	c := pbCatalog.NewCatalogServiceClient(conn)
	return &CatalogClient{
		client: c,
	}, nil
}

func (cc *CatalogClient) GetContent(
	ctx context.Context,
	contentID string,
) (*dto.Content, error) {

	resp, err := cc.client.GetContent(
		ctx,
		&pbCatalog.GetContentRequest{Id: contentID},
	)
	if err != nil {
		return nil, err
	}

	c := resp.GetContent()

	return &dto.Content{
		ID:             c.GetId(),
		ExternalSource: c.GetExternalSource(),
		ExternalID:     c.GetExternalId(),
		Title:          c.GetTitle(),
		Description:    c.GetDescription(),
		PosterURL:      c.GetPosterUrl(),
		ReleaseDate:    c.GetReleaseDate(),
	}, nil
}
