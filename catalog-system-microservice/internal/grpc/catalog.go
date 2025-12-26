package catalog

import (
	"catalog-microservice/internal/domain/models"
	catalog1 "catalog-microservice/internal/pb/catalog"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Catalog interface {
	CreateContent(ctx context.Context, context models.Content) (string, error)
	GetContent(ctx context.Context, id string) (models.Content, error)
	UpdateContent(ctx context.Context, id string, content models.Content) (models.Content, error)
	FindContentByExternal(ctx context.Context, externalID string) (models.Content, error)
	GetMovieDetails(ctx context.Context, id string) (models.MovieDetails, error)
	GetAnimeDetails(ctx context.Context, id string) (models.AnimeDetails, error)
	GetGameDetails(ctx context.Context, id string) (models.GameDetails, error)
	GetSeriesDetails(ctx context.Context, id string) (models.SeriesDetails, error)
	GetBookDetails(ctx context.Context, id string) (models.BookDetails, error)
}

type serverAPI struct {
	catalog1.UnimplementedCatalogServiceServer
	catalog Catalog
}

func Register(gRPC *grpc.Server, catalog Catalog) {
	catalog1.RegisterCatalogServiceServer(gRPC, &serverAPI{catalog: catalog})
}

func (s *serverAPI) CreateContent(ctx context.Context, req *catalog1.CreateContentRequest) (*catalog1.CreateContentResponse, error) {
	content := &models.Content{
		ID:             req.Content.GetId(),
		Type:           req.Content.GetType().String(),
		ExternalSource: req.Content.GetExternalSource(),
		ExternalID:     req.Content.GetExternalId(),
		Title:          req.Content.GetTitle(),
		Description:    req.Content.GetDescription(),
		PosterURL:      req.Content.GetPosterUrl(),
		ReleaseDate:    req.Content.GetReleaseDate(),
		CreatedAt:      req.Content.GetCreatedAt(),
		UpdatedAt:      req.Content.GetUpdatedAt(),
	}

	//TODO: validate content
	id, err := s.catalog.CreateContent(ctx, *content)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	content.ID = id

	//TODO: return created content
	return &catalog1.CreateContentResponse{
		Content: &catalog1.Content{},
	}, nil
}

func (s *serverAPI) GetContent(ctx context.Context, req *catalog1.GetContentRequest) (*catalog1.GetContentResponse, error) {
	content, err := s.catalog.GetContent(ctx, req.GetId())
	if err != nil {
		if req.GetId() == "" {
			return nil, status.Error(codes.InvalidArgument, "wrong argument")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog1.GetContentResponse{
		Content: &catalog1.Content{
			Id: content.ID,
			//Type:           content.Type,
			ExternalSource: content.ExternalSource,
			ExternalId:     content.ExternalID,
			Title:          content.Title,
			Description:    content.Description,
			PosterUrl:      content.PosterURL,
			ReleaseDate:    content.ReleaseDate,
			CreatedAt:      content.CreatedAt,
			UpdatedAt:      content.UpdatedAt,
		},
	}, nil
}

func (s *serverAPI) UpdateContent(ctx context.Context, req *catalog1.UpdateContentRequest) (*catalog1.UpdateContentResponse, error) {
	content := &models.Content{
		ID:   req.Content.GetId(),
		Type: req.Content.GetType().String(),
	}
	newContent, err := s.catalog.UpdateContent(ctx, req.Content.GetId(), *content)
	if err != nil {
		if req.Content.GetId() == "" {
			return nil, status.Error(codes.InvalidArgument, "wrong argument")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &catalog1.UpdateContentResponse{
		Content: &catalog1.Content{
			Id: newContent.ID,
		},
	}, nil
}

func (s *serverAPI) FindContentByExternal(ctx context.Context, req *catalog1.FindContentByExternalRequest) (*catalog1.FindContentByExternalResponse, error) {
	if req.GetExternalId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong argument")
	}

	content, err := s.catalog.FindContentByExternal(ctx, req.GetExternalId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog1.FindContentByExternalResponse{
		Content: &catalog1.Content{
			Id: content.ID,
		},
	}, nil
}

func (s *serverAPI) GetMovieDetails(ctx context.Context, req *catalog1.GetMovieDetailsRequest) (*catalog1.GetMovieDetailsResponse, error) {

	if req.GetContentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong argument")
	}

	movieDetails, err := s.catalog.GetMovieDetails(ctx, req.GetContentId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog1.GetMovieDetailsResponse{
		Details: &catalog1.MovieDetails{
			ContentId: movieDetails.ContentID,
		},
	}, nil
}

func (s *serverAPI) GetAnimeDetails(ctx context.Context, req *catalog1.GetAnimeDetailsRequest) (*catalog1.GetAnimeDetailsResponse, error) {

	if req.GetContentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong argument")
	}

	animeDetails, err := s.catalog.GetAnimeDetails(ctx, req.GetContentId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog1.GetAnimeDetailsResponse{
		Details: &catalog1.AnimeDetails{
			ContentId: animeDetails.ContentID,
		},
	}, nil
}

func (s *serverAPI) GetGameDetails(ctx context.Context, req *catalog1.GetGameDetailsRequest) (*catalog1.GetGameDetailsResponse, error) {

	if req.GetContentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong argument")
	}

	gameDetails, err := s.catalog.GetGameDetails(ctx, req.GetContentId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog1.GetGameDetailsResponse{
		Details: &catalog1.GameDetails{
			ContentId: gameDetails.ContentID,
		},
	}, nil
}

func (s *serverAPI) GetBookDetails(ctx context.Context, req *catalog1.GetBookDetailsRequest) (*catalog1.GetBookDetailsResponse, error) {

	if req.GetContentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong argument")
	}

	bookDetails, err := s.catalog.GetBookDetails(ctx, req.GetContentId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog1.GetBookDetailsResponse{
		Details: &catalog1.BookDetails{
			ContentId: bookDetails.ContentID,
		},
	}, nil
}

func (s *serverAPI) GetSeriesDetails(ctx context.Context, req *catalog1.GetSeriesDetailsRequest) (*catalog1.GetSeriesDetailsResponse, error) {

	if req.GetContentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong argument")
	}

	seriesDetails, err := s.catalog.GetSeriesDetails(ctx, req.GetContentId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog1.GetSeriesDetailsResponse{
		Details: &catalog1.SeriesDetails{
			ContentId: seriesDetails.ContentID,
		},
	}, nil
}
