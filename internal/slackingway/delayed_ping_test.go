package slackingway_test

import (
	"net/http"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestDelayedPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSlackingway := NewMockSlackingway(ctrl)
	mockSlackingway.EXPECT().NewResponse(gomock.Any()).Return(&http.Request{}, nil)
	mockSlackingway.EXPECT().SendResponse(gomock.Any()).Return(nil)

	slackingway.PingDelay = 1 * time.Millisecond

	expected := slack.Msg{Text: "Delayed Pong!", ReplaceOriginal: true}
	actual, err := slackingway.DelayedPing(mockSlackingway)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
