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

func TestCreateMenuModal(t *testing.T) {
	modal := slackingway.CreateMenuModal()
	assert.NotNil(t, modal)
}
