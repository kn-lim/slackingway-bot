package slackingway_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

func TestNewSlackingway(t *testing.T) {
	body := &slackingway.SlackRequestBody{
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

	actual := slackingway.NewSlackingway(body)

	assert.NotNil(t, actual)
	assert.Equal(t, body, actual.SlackRequestBody)
	assert.NotNil(t, actual.HTTPClient)
}

func TestNewResponse(t *testing.T) {
	message := slack.Msg{Text: "TestNewResponse"}

	s := slackingway.NewSlackingway(&slackingway.SlackRequestBody{})

	actual, err := s.NewResponse(message)

	assert.Nil(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, "POST", actual.Method)

	// Read the body from the actual request and compare it to the expected message
	var actualMessage slack.Msg
	bodyBytes, err := io.ReadAll(actual.Body) // Read the request body
	assert.Nil(t, err)

	// Unmarshal the body into a slack.Msg object
	err = json.Unmarshal(bodyBytes, &actualMessage)
	assert.Nil(t, err)

	// Assert that the actual message in the body matches the expected message
	assert.Equal(t, message, actualMessage)
}

func TestSendResponse(t *testing.T) {
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
