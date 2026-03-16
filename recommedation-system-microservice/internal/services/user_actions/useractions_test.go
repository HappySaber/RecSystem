package useractions_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"rec-system-microservice/internal/domain/models"
	useractions "rec-system-microservice/internal/services/user_actions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTracker struct{ mock.Mock }

func (m *mockTracker) SaveTrackAction(ctx context.Context, event models.UserActionEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func newTestUserActions() (*useractions.UserActionTracker, *mockTracker) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	tracker := new(mockTracker)
	svc := useractions.New(log, tracker)
	return svc, tracker
}

func TestTrackUserAction_Success(t *testing.T) {
	svc, tracker := newTestUserActions()

	event := models.UserActionEvent{
		UserID:    "user-1",
		ContentID: "content-1",
		Action:    models.ActionLike,
	}

	tracker.On("SaveTrackAction", mock.Anything, event).
		Return(nil).Once()

	err := svc.TrackUserAction(context.Background(), event)

	require.NoError(t, err)
	tracker.AssertExpectations(t)
}

func TestTrackUserAction_AllActionTypes(t *testing.T) {
	actions := []models.ActionType{
		models.ActionView,
		models.ActionLike,
		models.ActionDislike,
		models.ActionRate,
		models.ActionFavorite,
	}

	for _, action := range actions {
		t.Run(string(action), func(t *testing.T) {
			svc, tracker := newTestUserActions()

			event := models.UserActionEvent{
				UserID:    "user-1",
				ContentID: "content-1",
				Action:    action,
			}

			tracker.On("SaveTrackAction", mock.Anything, event).
				Return(nil).Once()

			err := svc.TrackUserAction(context.Background(), event)

			require.NoError(t, err)
			tracker.AssertExpectations(t)
		})
	}
}

func TestTrackUserAction_Error(t *testing.T) {
	svc, tracker := newTestUserActions()

	event := models.UserActionEvent{
		UserID:    "user-1",
		ContentID: "content-1",
		Action:    models.ActionView,
	}

	tracker.On("SaveTrackAction", mock.Anything, event).
		Return(errors.New("storage unavailable")).Once()

	err := svc.TrackUserAction(context.Background(), event)

	require.Error(t, err)
	assert.ErrorContains(t, err, "storage unavailable")
	tracker.AssertExpectations(t)
}
