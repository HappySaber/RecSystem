package client

import (
	"api-gateway/internal/dto"
	pbCatalog "api-gateway/internal/pbs/catalog"
	"context"
	"encoding/json"
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
		Type:           contentTypeToString(c.GetType()),
		ExternalSource: c.GetExternalSource(),
		ExternalID:     c.GetExternalId(),
		Title:          c.GetTitle(),
		Description:    c.GetDescription(),
		PosterURL:      c.GetPosterUrl(),
		ReleaseDate:    c.GetReleaseDate(),
	}, nil
}

func contentTypeToString(t pbCatalog.ContentType) string {
	switch t {
	case pbCatalog.ContentType_MOVIE:
		return "movie"
	case pbCatalog.ContentType_SERIES:
		return "series"
	case pbCatalog.ContentType_ANIME:
		return "anime"
	case pbCatalog.ContentType_GAME:
		return "game"
	case pbCatalog.ContentType_BOOK:
		return "book"
	default:
		return ""
	}
}

func parseStringSliceJSON(raw string) []string {
	if raw == "" || raw == "null" {
		return nil
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil
	}
	return items
}

func (cc *CatalogClient) GetContentByIDs(
	ctx context.Context,
	ids []string,
) ([]dto.ContentShort, error) {

	resp, err := cc.client.GetContentByIDs(
		ctx,
		&pbCatalog.GetContentByIDsRequest{Ids: ids},
	)
	if err != nil {
		return nil, err
	}

	out := make([]dto.ContentShort, 0, len(resp.GetContents()))

	for _, c := range resp.GetContents() {
		out = append(out, dto.ContentShort{
			ID: c.GetId(),
			//Type:  c.GetType(),
			Title: c.GetTitle(),
		})
	}

	return out, nil
}

func (cc *CatalogClient) FindContentByExternal(
	ctx context.Context,
	externalID string,
	externalSource string,
) (*dto.Content, error) {

	resp, err := cc.client.FindContentByExternal(
		ctx,
		&pbCatalog.FindContentByExternalRequest{
			ExternalId:     externalID,
			ExternalSource: externalSource,
		},
	)
	if err != nil {
		return nil, err
	}

	c := resp.GetContent()

	return &dto.Content{
		ID: c.GetId(),
		//Type:           c.GetType(),
		ExternalSource: c.GetExternalSource(),
		ExternalID:     c.GetExternalId(),
		Title:          c.GetTitle(),
		Description:    c.GetDescription(),
		PosterURL:      c.GetPosterUrl(),
		ReleaseDate:    c.GetReleaseDate(),
	}, nil
}

func (cc *CatalogClient) GetMovieDetails(
	ctx context.Context,
	id string,
) (*dto.MovieDetails, error) {

	resp, err := cc.client.GetMovieDetails(
		ctx,
		&pbCatalog.GetMovieDetailsRequest{ContentId: id},
	)
	if err != nil {
		return nil, err
	}

	m := resp.GetDetails()
	var runtime *int32
	if r := m.GetRuntime(); r > 0 {
		runtime = &r
	}

	return &dto.MovieDetails{
		ContentID:     m.GetContentId(),
		TmdbID:        m.GetTmdbId(),
		OriginalTitle: m.GetOriginalTitle(),
		Runtime:       runtime,
		Tagline:       m.GetTagline(),
		Status:        m.GetStatus(),
		Budget:        m.GetBudget(),
		Revenue:       m.GetRevenue(),
		Language:      m.GetLanguage(),
		Genres:        parseStringSliceJSON(m.GetGenresJson()),
		Cast:          parseStringSliceJSON(m.GetCastJson()),
	}, nil
}

func (cc *CatalogClient) GetAllMovieDetails(
	ctx context.Context,
) ([]dto.MovieDetails, error) {

	resp, err := cc.client.GetAllMovieDetails(
		ctx,
		&pbCatalog.GetAllMovieDetailsRequest{},
	)
	if err != nil {
		return nil, err
	}

	out := make([]dto.MovieDetails, 0, len(resp.GetDetails()))

	for _, m := range resp.GetDetails() {
		out = append(out, dto.MovieDetails{
			ContentID:     m.GetContentId(),
			TmdbID:        m.GetTmdbId(),
			OriginalTitle: m.GetOriginalTitle(),
			Language:      m.GetLanguage(),
		})
	}

	return out, nil
}

func (cc *CatalogClient) GetSeriesDetails(
	ctx context.Context,
	id string,
) (*dto.SeriesDetails, error) {

	resp, err := cc.client.GetSeriesDetails(
		ctx,
		&pbCatalog.GetSeriesDetailsRequest{ContentId: id},
	)
	if err != nil {
		return nil, err
	}

	s := resp.GetDetails()

	return &dto.SeriesDetails{
		ContentID:        s.GetContentId(),
		TmdbID:           s.GetTmdbId(),
		OriginalName:     s.GetOriginalName(),
		Status:           s.GetStatus(),
		NumberOfSeasons:  s.GetNumberOfSeasons(),
		NumberOfEpisodes: s.GetNumberOfEpisodes(),
		Language:         s.GetLanguage(),
	}, nil
}

func (cc *CatalogClient) GetAllSeriesDetails(
	ctx context.Context,
	id string,
) ([]dto.SeriesDetails, error) {

	resp, err := cc.client.GetAllSeriesDetails(
		ctx,
		&pbCatalog.GetAllSeriesDetailsRequest{},
	)
	if err != nil {
		return nil, err
	}

	out := make([]dto.SeriesDetails, 0, len(resp.GetDetails()))

	for _, m := range resp.GetDetails() {
		out = append(out, dto.SeriesDetails{
			ContentID:    m.GetContentId(),
			TmdbID:       m.GetTmdbId(),
			OriginalName: m.GetOriginalName(),
			Language:     m.GetLanguage(),
		})
	}

	return out, nil
}
