package slackingway_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

// TestDelayedPing tests the DelayedPing function
func TestDelayedPing(t *testing.T) {
	// Create a new mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Error on NewResponse", func(t *testing.T) {
		mockSlackingway := NewMockSlackingway(ctrl)
		mockSlackingway.EXPECT().NewResponse(gomock.Any()).Return(&http.Request{}, errors.New("error!!!"))

		_, err := slackingway.DelayedPing(mockSlackingway)
		assert.NotNil(t, err)
	})

	t.Run("Error on SendResponse", func(t *testing.T) {
		mockSlackingway := NewMockSlackingway(ctrl)
		mockSlackingway.EXPECT().NewResponse(gomock.Any()).Return(&http.Request{}, nil)
		mockSlackingway.EXPECT().SendResponse(gomock.Any()).Return(errors.New("error!!!"))

		_, err := slackingway.DelayedPing(mockSlackingway)
		assert.NotNil(t, err)
	})

	t.Run("Successful DelayedPing", func(t *testing.T) {
		mockSlackingway := NewMockSlackingway(ctrl)
		mockSlackingway.EXPECT().NewResponse(gomock.Any()).Return(&http.Request{}, nil)
		mockSlackingway.EXPECT().SendResponse(gomock.Any()).Return(nil)

		// Set PingDelay to a minimal value for testing
		slackingway.PingDelay = 1 * time.Millisecond

		expected := slack.Msg{Text: "Delayed Pong!", ReplaceOriginal: true}
		actual, err := slackingway.DelayedPing(mockSlackingway)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
