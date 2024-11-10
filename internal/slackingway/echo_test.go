package slackingway_test

import (
	"errors"
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

	t.Run("Error on OpenView", func(t *testing.T) {
		MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
		MockSlackAPIClient.EXPECT().OpenView(gomock.Any(), gomock.Any()).Return(&slack.ViewResponse{}, errors.New("error!!!"))

		s := &slackingway.SlackingwayWrapper{
			APIClient:        MockSlackAPIClient,
			SlackRequestBody: &slackingway.SlackRequestBody{TriggerID: "definitely_a_valid_trigger_id"},
		}
		err := slackingway.Echo(s)

		assert.NotNil(t, err)
	})

	t.Run("Successful Echo", func(t *testing.T) {
		MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
		MockSlackAPIClient.EXPECT().OpenView(gomock.Any(), gomock.Any()).Return(&slack.ViewResponse{}, nil)

		s := &slackingway.SlackingwayWrapper{
			APIClient:        MockSlackAPIClient,
			SlackRequestBody: &slackingway.SlackRequestBody{TriggerID: "definitely_a_valid_trigger_id"},
		}
		err := slackingway.Echo(s)

		assert.Nil(t, err)
	})
}

func TestReceivedEcho(t *testing.T) {
	// Create a new mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testView := slack.View{
		State: &slack.ViewState{
			Values: map[string]map[string]slack.BlockAction{},
		},
	}

	t.Run("Error on GetUserInfo", func(t *testing.T) {
		MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
		MockSlackAPIClient.EXPECT().GetUserInfo(gomock.Any()).Return(&slack.User{}, errors.New("error!!!"))

		s := &slackingway.SlackingwayWrapper{
			APIClient: MockSlackAPIClient,
			SlackRequestBody: &slackingway.SlackRequestBody{
				UserID: "definitely_a_invalid_user_id",
				View:   testView,
			},
		}
		msg, err := slackingway.ReceivedEcho(s)

		assert.NotNil(t, err)
		assert.Empty(t, msg)
	})

	t.Run("Successful ReceivedEcho", func(t *testing.T) {
		MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
		MockSlackAPIClient.EXPECT().GetUserInfo(gomock.Any()).Return(&slack.User{}, nil)

		s := &slackingway.SlackingwayWrapper{
			APIClient: MockSlackAPIClient,
			SlackRequestBody: &slackingway.SlackRequestBody{
				UserID: "definitely_a_valid_user_id",
				View:   testView,
			},
		}
		msg, err := slackingway.ReceivedEcho(s)

		assert.Nil(t, err)
		assert.NotEmpty(t, msg)
	})
}

func TestGenerateEchoModal(t *testing.T) {
	got := slackingway.CreateEchoModal()

	// Run test
	assert.NotNil(t, got)
}
