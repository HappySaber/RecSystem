package tests

import (
	"testing"

	catalogpb "catalog-microservice/internal/pb/catalog"
	"catalog-microservice/tests/suite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GetContent

func TestGetContent_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	// вставляем тестовые данные
	id := st.InsertContent(t, "movie", "tmdb", "550", "Fight Club")

	resp, err := st.CatalogClient.GetContent(ctx, &catalogpb.GetContentRequest{
		Id: id,
	})

	require.NoError(t, err)
	assert.Equal(t, id, resp.GetContent().GetId())
	assert.Equal(t, "Fight Club", resp.GetContent().GetTitle())
	assert.Equal(t, "tmdb", resp.GetContent().GetExternalSource())
	assert.Equal(t, "550", resp.GetContent().GetExternalId())
}

func TestGetContent_EmptyID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.GetContent(ctx, &catalogpb.GetContentRequest{
		Id: "",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong argument")
}

func TestGetContent_NotFound(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.GetContent(ctx, &catalogpb.GetContentRequest{
		Id: "00000000-0000-0000-0000-000000000000",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}

// GetContentByIDs

func TestGetContentByIDs_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	id1 := st.InsertContent(t, "movie", "tmdb", "551", "Inception")
	id2 := st.InsertContent(t, "anime", "anilist", "1", "Naruto")

	resp, err := st.CatalogClient.GetContentByIDs(ctx, &catalogpb.GetContentByIDsRequest{
		Ids: []string{id1, id2},
	})

	require.NoError(t, err)
	assert.Len(t, resp.GetContents(), 2)
}

func TestGetContentByIDs_EmptyList(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.GetContentByIDs(ctx, &catalogpb.GetContentByIDsRequest{
		Ids: []string{},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong argument")
}

func TestGetContentByIDs_NonExistentIDs(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.CatalogClient.GetContentByIDs(ctx, &catalogpb.GetContentByIDsRequest{
		Ids: []string{
			"00000000-0000-0000-0000-000000000001",
			"00000000-0000-0000-0000-000000000002",
		},
	})

	// не ошибка — просто пустой список
	require.NoError(t, err)
	assert.Empty(t, resp.GetContents())
}

// FindContentByExternal

func TestFindContentByExternal_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	st.InsertContent(t, "movie", "tmdb", "999", "The Matrix")

	resp, err := st.CatalogClient.FindContentByExternal(ctx, &catalogpb.FindContentByExternalRequest{
		ExternalId:     "999",
		ExternalSource: "tmdb",
	})

	require.NoError(t, err)
	assert.Equal(t, "The Matrix", resp.GetContent().GetTitle())
	assert.Equal(t, "999", resp.GetContent().GetExternalId())
}

func TestFindContentByExternal_EmptyID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.FindContentByExternal(ctx, &catalogpb.FindContentByExternalRequest{
		ExternalId:     "",
		ExternalSource: "tmdb",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong argument")
}

func TestFindContentByExternal_NotFound(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.FindContentByExternal(ctx, &catalogpb.FindContentByExternalRequest{
		ExternalId:     "nonexistent-99999",
		ExternalSource: "tmdb",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}

// GetMovieDetails

func TestGetMovieDetails_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	id := st.InsertContent(t, "movie", "tmdb", "552", "Interstellar")
	st.InsertMovieDetails(t, id, 157336)

	resp, err := st.CatalogClient.GetMovieDetails(ctx, &catalogpb.GetMovieDetailsRequest{
		ContentId: id,
	})

	require.NoError(t, err)
	assert.Equal(t, id, resp.GetDetails().GetContentId())
}

func TestGetMovieDetails_EmptyID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.GetMovieDetails(ctx, &catalogpb.GetMovieDetailsRequest{
		ContentId: "",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong argument")
}

func TestGetMovieDetails_NotFound(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.GetMovieDetails(ctx, &catalogpb.GetMovieDetailsRequest{
		ContentId: "00000000-0000-0000-0000-000000000000",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}

// GetAnimeDetails

func TestGetAnimeDetails_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	id := st.InsertContent(t, "anime", "anilist", "20", "Naruto")
	st.InsertAnimeDetails(t, id, 20)

	resp, err := st.CatalogClient.GetAnimeDetails(ctx, &catalogpb.GetAnimeDetailsRequest{
		ContentId: id,
	})

	require.NoError(t, err)
	assert.Equal(t, id, resp.GetDetails().GetContentId())
}

func TestGetAnimeDetails_EmptyID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.GetAnimeDetails(ctx, &catalogpb.GetAnimeDetailsRequest{
		ContentId: "",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong argument")
}

func TestGetAnimeDetails_NotFound(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.GetAnimeDetails(ctx, &catalogpb.GetAnimeDetailsRequest{
		ContentId: "00000000-0000-0000-0000-000000000000",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}

// GetGameDetails

// func TestGetGameDetails_HappyPath(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	id := st.InsertContent(t, "game", "igdb", "1942", "The Witcher 3")
// 	st.InsertGameDetails(t, id, 1942)

// 	resp, err := st.CatalogClient.GetGameDetails(ctx, &catalogpb.GetGameDetailsRequest{
// 		ContentId: id,
// 	})

// 	require.NoError(t, err)
// 	assert.Equal(t, id, resp.GetDetails().GetContentId())
// }

// func TestGetGameDetails_EmptyID(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	_, err := st.CatalogClient.GetGameDetails(ctx, &catalogpb.GetGameDetailsRequest{
// 		ContentId: "",
// 	})

// 	require.Error(t, err)
// 	assert.ErrorContains(t, err, "wrong argument")
// }

// func TestGetGameDetails_NotFound(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	_, err := st.CatalogClient.GetGameDetails(ctx, &catalogpb.GetGameDetailsRequest{
// 		ContentId: "00000000-0000-0000-0000-000000000000",
// 	})

// 	require.Error(t, err)
// 	assert.ErrorContains(t, err, "not found")
// }

// GetSeriesDetails

func TestGetSeriesDetails_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	id := st.InsertContent(t, "series", "tmdb", "1396", "Breaking Bad")
	st.InsertSeriesDetails(t, id, 1396)

	resp, err := st.CatalogClient.GetSeriesDetails(ctx, &catalogpb.GetSeriesDetailsRequest{
		ContentId: id,
	})

	require.NoError(t, err)
	assert.Equal(t, id, resp.GetDetails().GetContentId())
}

func TestGetSeriesDetails_EmptyID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.GetSeriesDetails(ctx, &catalogpb.GetSeriesDetailsRequest{
		ContentId: "",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong argument")
}

func TestGetSeriesDetails_NotFound(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.CatalogClient.GetSeriesDetails(ctx, &catalogpb.GetSeriesDetailsRequest{
		ContentId: "00000000-0000-0000-0000-000000000000",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}

// GetBookDetails

// func TestGetBookDetails_HappyPath(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	id := st.InsertContent(t, "book", "openlibrary", "OL7353617M", "Dune")
// 	st.InsertBookDetails(t, id)

// 	resp, err := st.CatalogClient.GetBookDetails(ctx, &catalogpb.GetBookDetailsRequest{
// 		ContentId: id,
// 	})

// 	require.NoError(t, err)
// 	assert.Equal(t, id, resp.GetDetails().GetContentId())
// }

// func TestGetBookDetails_EmptyID(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	_, err := st.CatalogClient.GetBookDetails(ctx, &catalogpb.GetBookDetailsRequest{
// 		ContentId: "",
// 	})

// 	require.Error(t, err)
// 	assert.ErrorContains(t, err, "wrong argument")
// }

// GetAll методы

func TestGetAllMovieDetails_EmptyDB(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.CatalogClient.GetAllMovieDetails(ctx, &catalogpb.GetAllMovieDetailsRequest{})

	require.NoError(t, err)
	assert.Empty(t, resp.GetDetails())
}

func TestGetAllMovieDetails_WithData(t *testing.T) {
	ctx, st := suite.New(t)

	id1 := st.InsertContent(t, "movie", "tmdb", "553", "Avatar")
	id2 := st.InsertContent(t, "movie", "tmdb", "554", "Titanic")
	st.InsertMovieDetails(t, id1, 19995)
	st.InsertMovieDetails(t, id2, 597)

	resp, err := st.CatalogClient.GetAllMovieDetails(ctx, &catalogpb.GetAllMovieDetailsRequest{})

	require.NoError(t, err)
	assert.Len(t, resp.GetDetails(), 2)
}

func TestGetAllAnimeDetails_EmptyDB(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.CatalogClient.GetAllAnimeDetails(ctx, &catalogpb.GetAllAnimeDetailsRequest{})

	require.NoError(t, err)
	assert.Empty(t, resp.GetDetails())
}

func TestGetAllAnimeDetails_WithData(t *testing.T) {
	ctx, st := suite.New(t)

	id1 := st.InsertContent(t, "anime", "anilist", "21", "One Piece")
	id2 := st.InsertContent(t, "anime", "anilist", "16498", "Attack on Titan")
	st.InsertAnimeDetails(t, id1, 21)
	st.InsertAnimeDetails(t, id2, 16498)

	resp, err := st.CatalogClient.GetAllAnimeDetails(ctx, &catalogpb.GetAllAnimeDetailsRequest{})

	require.NoError(t, err)
	assert.Len(t, resp.GetDetails(), 2)
}

// func TestGetAllGameDetails_EmptyDB(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	resp, err := st.CatalogClient.GetAllGameDetails(ctx, &catalogpb.GetAllGameDetailsRequest{})

// 	require.NoError(t, err)
// 	assert.Empty(t, resp.GetDetails())
// }

func TestGetAllSeriesDetails_EmptyDB(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.CatalogClient.GetAllSeriesDetails(ctx, &catalogpb.GetAllSeriesDetailsRequest{})

	require.NoError(t, err)
	assert.Empty(t, resp.GetDetails())
}

// func TestGetAllBookDetails_EmptyDB(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	resp, err := st.CatalogClient.GetAllBookDetails(ctx, &catalogpb.GetAllBookDetailsRequest{})

// 	require.NoError(t, err)
// 	assert.Empty(t, resp.GetDetails())
// }
