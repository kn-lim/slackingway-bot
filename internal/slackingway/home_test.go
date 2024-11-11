package slackingway_test

import (
	"errors"
	"testing"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHomeTab(t *testing.T) {
	// Create a new mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserID := "U12345678"

	t.Run("Error on GetUserInfo", func(t *testing.T) {
		MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
		MockSlackAPIClient.EXPECT().GetUserInfo(gomock.Any()).Return(&slack.User{}, errors.New("error!!!"))

		s := &slackingway.SlackingwayWrapper{
			APIClient: MockSlackAPIClient,
			SlackRequestBody: &slackingway.SlackRequestBody{
				UserID: testUserID,
			},
		}
		err := slackingway.HomeTab(s, testUserID)

		assert.NotNil(t, err)
	})

	t.Run("Error on PublishView", func(t *testing.T) {
		MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
		MockSlackAPIClient.EXPECT().GetUserInfo(gomock.Any()).Return(&slack.User{}, nil)
		MockSlackAPIClient.EXPECT().PublishView(gomock.Any(), gomock.Any(), gomock.Any()).Return(&slack.ViewResponse{}, errors.New("error!!!"))

		s := &slackingway.SlackingwayWrapper{
			APIClient: MockSlackAPIClient,
			SlackRequestBody: &slackingway.SlackRequestBody{
				UserID: testUserID,
			},
		}
		err := slackingway.HomeTab(s, testUserID)

		assert.NotNil(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
		MockSlackAPIClient.EXPECT().GetUserInfo(gomock.Any()).Return(&slack.User{}, nil)
		MockSlackAPIClient.EXPECT().PublishView(gomock.Any(), gomock.Any(), gomock.Any()).Return(&slack.ViewResponse{}, nil)

		s := &slackingway.SlackingwayWrapper{
			APIClient: MockSlackAPIClient,
			SlackRequestBody: &slackingway.SlackRequestBody{
				UserID: testUserID,
			},
		}
		err := slackingway.HomeTab(s, testUserID)

		assert.Nil(t, err)
	})
}
