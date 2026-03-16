package tests

import (
	"testing"

	recpb "rec-system-microservice/internal/pb/recommendation"
	"rec-system-microservice/tests/suite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// uuid константы для тестов — фиксированные чтобы легко читать
const (
	userID     = "00000000-0000-0000-0000-000000000001"
	contentID1 = "00000000-0000-0000-0000-000000000010"
	contentID2 = "00000000-0000-0000-0000-000000000011"
	contentID3 = "00000000-0000-0000-0000-000000000012"
)

// TrackUserAction

func TestTrackUserAction_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.Client.TrackUserAction(ctx, &recpb.TrackUserActionRequest{
		UserId:    userID,
		ContentId: contentID1,
		Action:    recpb.UserAction_LIKE,
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTrackUserAction_AllActionTypes(t *testing.T) {
	actions := []recpb.UserAction{
		recpb.UserAction_VIEW,
		recpb.UserAction_LIKE,
		recpb.UserAction_DISLIKE,
		recpb.UserAction_RATE,
		recpb.UserAction_ADD_TO_FAVORITES,
	}

	for _, action := range actions {
		action := action
		t.Run(action.String(), func(t *testing.T) {
			ctx, st := suite.New(t)

			_, err := st.Client.TrackUserAction(ctx, &recpb.TrackUserActionRequest{
				UserId:    userID,
				ContentId: contentID1,
				Action:    action,
			})

			require.NoError(t, err)
		})
	}
}

func TestTrackUserAction_EmptyUserID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.Client.TrackUserAction(ctx, &recpb.TrackUserActionRequest{
		UserId:    "",
		ContentId: contentID1,
		Action:    recpb.UserAction_LIKE,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong arguments")
}

func TestTrackUserAction_EmptyContentID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.Client.TrackUserAction(ctx, &recpb.TrackUserActionRequest{
		UserId:    userID,
		ContentId: "",
		Action:    recpb.UserAction_LIKE,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong arguments")
}

// GetRecommendations

func TestGetRecommendations_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	// создаём жанры
	actionGenreID := st.InsertGenre(t, "Action")
	dramaGenreID := st.InsertGenre(t, "Drama")

	// contentID1 — пользователь уже смотрел, не должен попасть в рекомендации
	st.InsertContentGenre(t, contentID1, actionGenreID)
	st.InsertUserAction(t, userID, contentID1, "LIKE")

	// contentID2 — тот же жанр что смотрел, должен попасть в рекомендации
	st.InsertContentGenre(t, contentID2, actionGenreID)

	// contentID3 — другой жанр, меньший приоритет
	st.InsertContentGenre(t, contentID3, dramaGenreID)

	resp, err := st.Client.GetRecommendations(ctx, &recpb.GetRecommendationsRequest{
		UserId: userID,
		Limit:  10,
	})

	require.NoError(t, err)
	// contentID1 не должен быть в результатах — пользователь уже взаимодействовал
	assert.NotContains(t, resp.GetContentIds(), contentID1)
	// contentID2 должен быть — тот же жанр что лайкнул пользователь
	assert.Contains(t, resp.GetContentIds(), contentID2)
}

func TestGetRecommendations_EmptyUserID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.Client.GetRecommendations(ctx, &recpb.GetRecommendationsRequest{
		UserId: "",
		Limit:  10,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong arguments")
}

func TestGetRecommendations_ZeroLimit(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.Client.GetRecommendations(ctx, &recpb.GetRecommendationsRequest{
		UserId: userID,
		Limit:  0,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong arguments")
}

func TestGetRecommendations_NoHistory(t *testing.T) {
	ctx, st := suite.New(t)

	// пользователь без истории — пустой результат это норма
	resp, err := st.Client.GetRecommendations(ctx, &recpb.GetRecommendationsRequest{
		UserId: "00000000-0000-0000-0000-000000000099",
		Limit:  10,
	})

	require.NoError(t, err)
	assert.Empty(t, resp.GetContentIds())
}

// GetRecommendationsByGenres

func TestGetRecommendationsByGenres_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	genreID := st.InsertGenre(t, "Thriller")
	st.InsertContentGenre(t, contentID1, genreID)
	st.InsertContentGenre(t, contentID2, genreID)

	resp, err := st.Client.GetRecommendationsByGenres(ctx, &recpb.GetRecommendationsByGenresRequest{
		Genres: []string{"Thriller"},
		Limit:  10,
	})

	require.NoError(t, err)
	assert.Contains(t, resp.GetContentIds(), contentID1)
	assert.Contains(t, resp.GetContentIds(), contentID2)
}

func TestGetRecommendationsByGenres_EmptyGenres(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.Client.GetRecommendationsByGenres(ctx, &recpb.GetRecommendationsByGenresRequest{
		Genres: []string{},
		Limit:  10,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong arguments")
}

func TestGetRecommendationsByGenres_UnknownGenre(t *testing.T) {
	ctx, st := suite.New(t)

	// жанр которого нет в БД — пустой результат
	resp, err := st.Client.GetRecommendationsByGenres(ctx, &recpb.GetRecommendationsByGenresRequest{
		Genres: []string{"NonExistentGenre"},
		Limit:  10,
	})

	require.NoError(t, err)
	assert.Empty(t, resp.GetContentIds())
}

// GetSimilarContent

func TestGetSimilarContent_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	genreID := st.InsertGenre(t, "SciFi")

	// contentID1 и contentID2 имеют общий жанр
	st.InsertContentGenre(t, contentID1, genreID)
	st.InsertContentGenre(t, contentID2, genreID)

	resp, err := st.Client.GetSimilarContent(ctx, &recpb.GetSimilarContentRequest{
		ContentId: contentID1,
		Limit:     10,
	})

	require.NoError(t, err)
	// contentID2 похож на contentID1 по жанру
	assert.Contains(t, resp.GetContentIds(), contentID2)
	// сам контент не должен быть в результатах
	assert.NotContains(t, resp.GetContentIds(), contentID1)
}

func TestGetSimilarContent_EmptyContentID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.Client.GetSimilarContent(ctx, &recpb.GetSimilarContentRequest{
		ContentId: "",
		Limit:     10,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong arguments")
}

func TestGetSimilarContent_NoSimilar(t *testing.T) {
	ctx, st := suite.New(t)

	// контент без жанров — нет похожих
	resp, err := st.Client.GetSimilarContent(ctx, &recpb.GetSimilarContentRequest{
		ContentId: "00000000-0000-0000-0000-000000000099",
		Limit:     10,
	})

	require.NoError(t, err)
	assert.Empty(t, resp.GetContentIds())
}

// GetTrendingContent

func TestGetTrendingContent_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	// вставляем VIEW действие — trending считает просмотры за последние 24 часа
	st.InsertUserAction(t, userID, contentID1, "VIEW")
	st.InsertUserAction(t, userID, contentID2, "VIEW")

	resp, err := st.Client.GetTrendingContent(ctx, &recpb.GetTrendingContentRequest{
		Limit: 10,
	})

	require.NoError(t, err)
	assert.Contains(t, resp.GetContentIds(), contentID1)
	assert.Contains(t, resp.GetContentIds(), contentID2)
}

func TestGetTrendingContent_ZeroLimit(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.Client.GetTrendingContent(ctx, &recpb.GetTrendingContentRequest{
		Limit: 0,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong arguments")
}

func TestGetTrendingContent_EmptyDB(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.Client.GetTrendingContent(ctx, &recpb.GetTrendingContentRequest{
		Limit: 10,
	})

	require.NoError(t, err)
	assert.Empty(t, resp.GetContentIds())
}

// GetUserPreferences

func TestGetUserPreferences_EmptyUserID(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.Client.GetUserPreferences(ctx, &recpb.GetUserPreferencesRequest{
		UserId: "",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "wrong arguments")
}

func TestSetUserPreferences_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.Client.SetUserPreferences(ctx, &recpb.SetUserPreferencesRequest{
		UserId:          userID,
		PreferredGenres: []string{"Action", "Drama"},
	})

	require.NoError(t, err)
	assert.True(t, resp.GetUserPreferencesSetted())
}

func TestResetUserPreferences_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.Client.ResetUserPreferences(ctx, &recpb.ResetUserPreferencesRequest{
		UserId: userID,
	})

	require.NoError(t, err)
	assert.True(t, resp.GetUserPreferencesResetted())
}
