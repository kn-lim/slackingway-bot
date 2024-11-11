package slackingway_test

import (
	"errors"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

func TestMenu(t *testing.T) {
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
		err := slackingway.Menu(s)

		assert.NotNil(t, err)
	})

	t.Run("Successful Menu", func(t *testing.T) {
		MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
		MockSlackAPIClient.EXPECT().OpenView(gomock.Any(), gomock.Any()).Return(&slack.ViewResponse{}, nil)

		s := &slackingway.SlackingwayWrapper{
			APIClient:        MockSlackAPIClient,
			SlackRequestBody: &slackingway.SlackRequestBody{TriggerID: "definitely_a_valid_trigger_id"},
		}
		err := slackingway.Menu(s)

		assert.Nil(t, err)
	})
}

func TestReceivedMenu(t *testing.T) {
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
		blocks, err := slackingway.ReceivedMenu(s)

		assert.NotNil(t, err)
		assert.Empty(t, blocks)
	})

	t.Run("Successful ReceivedMenu", func(t *testing.T) {
		MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
		MockSlackAPIClient.EXPECT().GetUserInfo(gomock.Any()).Return(&slack.User{RealName: "DefinitelyA RealName"}, nil)

		s := &slackingway.SlackingwayWrapper{
			APIClient: MockSlackAPIClient,
			SlackRequestBody: &slackingway.SlackRequestBody{
				UserID: "definitely_a_valid_user_id",
				View:   testView,
			},
		}
		blocks, err := slackingway.ReceivedMenu(s)

		assert.Nil(t, err)
		assert.NotEmpty(t, blocks)
	})
}

func TestCreateMenuModal(t *testing.T) {
	modal := slackingway.CreateMenuModal()
	assert.NotNil(t, modal)
}
