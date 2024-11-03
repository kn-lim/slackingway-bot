package slackingway_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

// TestSlackRequestBody is a test SlackRequestBody for use in tests
var TestSlackRequestBody = &slackingway.SlackRequestBody{
	Timestamp:   "1234567890",
	Type:        "command",
	Challenge:   "challenge",
	Token:       "token",
	Command:     "/test-command",
	Text:        "testing123!",
	ResponseURL: "http://definitely-a-real-url.com/response",
	UserID:      "U12345",
	ChannelID:   "C12345",
	TeamID:      "T12345",
}

// TestNewSlackingway tests the NewSlackingway function
func TestNewSlackingway(t *testing.T) {
	// Run tests
	actual := slackingway.NewSlackingway(TestSlackRequestBody)

	assert.NotNil(t, actual)
	assert.Equal(t, TestSlackRequestBody, actual.SlackRequestBody)
	assert.NotNil(t, actual.HTTPClient)
	assert.NotNil(t, actual.APIClient)
}

// TestNewResponse tests the NewResponse function
func TestNewResponse(t *testing.T) {
	message := slack.Msg{Text: "TestNewResponse"}

	// Create a new SlackingwayWrapper instance
	s := slackingway.NewSlackingway(&slackingway.SlackRequestBody{})

	// Run tests
	actual, err := s.NewResponse(message)

	assert.Nil(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, "POST", actual.Method)

	// Read the body from the actual request and compare it to the expected message
	var actualMessage slack.Msg
	bodyBytes, err := io.ReadAll(actual.Body)
	assert.Nil(t, err)

	// Unmarshal the body into a slack.Msg object
	err = json.Unmarshal(bodyBytes, &actualMessage)
	assert.Nil(t, err)

	// Assert that the actual message in the body matches the expected message
	assert.Equal(t, message, actualMessage)
}

// TestSendResponse tests the SendResponse function
func TestSendResponse(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer mockServer.Close()

	// Create a SlackingwayWrapper instance with the mock server's URL
	s := &slackingway.SlackingwayWrapper{
		HTTPClient: mockServer.Client(),
	}

	// Create a new request
	req, err := http.NewRequest("POST", mockServer.URL, nil)
	assert.NoError(t, err)

	// Call SendResponse
	err = s.SendResponse(req)

	assert.Nil(t, err)
}

// TestWriteToHistory tests the WriteToHistory function
func TestWriteToHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create an empty SlackingwayWrapper instance
	s := &slackingway.SlackingwayWrapper{}

	assert.NotNil(t, s.WriteToHistory())

	mockSlackAPIClient := NewMockSlackAPIClient(ctrl)
	mockSlackAPIClient.EXPECT().GetUserInfo(gomock.Any()).Return(&slack.User{RealName: "DefinitelyA RealName"}, nil)
	mockSlackAPIClient.EXPECT().PostMessage(gomock.Any(), gomock.Any()).Return("messageID", "timestamp", nil)

	s = &slackingway.SlackingwayWrapper{
		APIClient:        mockSlackAPIClient,
		SlackRequestBody: TestSlackRequestBody,
	}

	assert.Nil(t, s.WriteToHistory())
}

// TestSendTextMessage tests the SendTextMessage function
func TestSendTextMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	message := "TestSendTextMessage"
	channelID := "C12345"

	// Create an empty SlackingwayWrapper instance
	s := &slackingway.SlackingwayWrapper{}

	assert.NotNil(t, s.SendTextMessage(message, channelID))

	mockSlackAPIClient := NewMockSlackAPIClient(ctrl)
	mockSlackAPIClient.EXPECT().PostMessage(gomock.Any(), gomock.Any()).Return("messageID", "timestamp", nil)

	s = &slackingway.SlackingwayWrapper{
		APIClient: mockSlackAPIClient,
	}

	assert.Nil(t, s.SendTextMessage(message, channelID))
}

// TestSendBlockMessage tests the SendTextMessage function
func TestSendBlockMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	blocks := []slack.Block{
		&slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "TestSendBlockMessage",
			},
		},
	}
	channelID := "C12345"

	// Create an empty SlackingwayWrapper instance
	s := &slackingway.SlackingwayWrapper{}

	assert.NotNil(t, s.SendBlockMessage(blocks, channelID))

	mockSlackAPIClient := NewMockSlackAPIClient(ctrl)
	mockSlackAPIClient.EXPECT().PostMessage(gomock.Any(), gomock.Any()).Return("messageID", "timestamp", nil)

	s = &slackingway.SlackingwayWrapper{
		APIClient: mockSlackAPIClient,
	}

	assert.Nil(t, s.SendBlockMessage(blocks, channelID))
}
