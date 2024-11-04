package slackingway_test

import (
	"testing"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEcho(t *testing.T) {
	// Create a new mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
	MockSlackAPIClient.EXPECT().OpenView(gomock.Any(), gomock.Any()).Return(&slack.ViewResponse{}, nil)

	s := &slackingway.SlackingwayWrapper{
		APIClient:        MockSlackAPIClient,
		SlackRequestBody: &slackingway.SlackRequestBody{TriggerID: "definitely_a_valid_trigger_id"},
	}
	err := slackingway.Echo(s)

	// Run test
	assert.Nil(t, err)
}

func TestGenerateEchoModal(t *testing.T) {
	got := slackingway.GenerateEchoModal()

	// Run test
	assert.NotNil(t, got)
}

func TestUpdateEchoModal(t *testing.T) {
	got := slackingway.UpdateEchoModal()

	// Run test
	assert.NotNil(t, got)
}
