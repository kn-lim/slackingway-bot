package slackingway_test

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

func TestParseTimestamp(t *testing.T) {
	var s slackingway.SlackRequestBody

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	assert.Nil(t, s.ParseTimestamp(timestamp))
}

func TestParsePayload(t *testing.T) {
	var s slackingway.SlackRequestBody

	t.Run("Invalid payload", func(t *testing.T) {
		invalidPayload := `invalid-payload`
		err := s.ParsePayload(invalidPayload)
		assert.NotNil(t, err)
	})

	t.Run("Valid payload", func(t *testing.T) {
		payload := `{
			"callback_id": "test_callback_id",
			"user": {
				"id": "U123456"
			},
			"channel_id": "C123456",
			"team_id": "T123456",
			"trigger_id": "123456.123456",
			"view": {
				"callback_id": "test_callback_id"
			}
		}`

		err := s.ParsePayload(payload)
		assert.Nil(t, err)
		assert.Equal(t, "test_callback_id", s.CallbackID)
		assert.Equal(t, "test_callback_id", s.Command)
		assert.Equal(t, "U123456", s.UserID)
		assert.Equal(t, "C123456", s.ChannelID)
		assert.Equal(t, "T123456", s.TeamID)
		assert.Equal(t, "123456.123456", s.TriggerID)
		assert.Equal(t, "test_callback_id", s.View.CallbackID)
	})
}

func TestParseSlashCommand(t *testing.T) {
	var s slackingway.SlackRequestBody

	t.Run("Valid values", func(t *testing.T) {
		requestData := url.Values{}
		requestData.Set("token", "test-token")
		requestData.Set("command", "/test-command")
		requestData.Set("text", "test-text")
		requestData.Set("response_url", "https://definitely-a-real-slack-response-url.com")
		requestData.Set("user_id", "U123456")
		requestData.Set("channel_id", "C123456")
		requestData.Set("team_id", "T123456")
		requestData.Set("callback_id", "test-callback-id")
		requestData.Set("trigger_id", "123456.123456")
		requestData.Set("view", `{"callback_id": "test-view-callback-id"}`)

		// Call the ParseSlashCommand method
		err := s.ParseSlashCommand(requestData)
		assert.Nil(t, err)

		// Verify that the fields are correctly set
		assert.Equal(t, "slash_command", s.Type)
		assert.Equal(t, "test-token", s.Token)
		assert.Equal(t, "/test-command", s.Command)
		assert.Equal(t, "test-text", s.Text)
		assert.Equal(t, "https://definitely-a-real-slack-response-url.com", s.ResponseURL)
		assert.Equal(t, "U123456", s.UserID)
		assert.Equal(t, "C123456", s.ChannelID)
		assert.Equal(t, "T123456", s.TeamID)
		assert.Equal(t, "test-callback-id", s.CallbackID)
		assert.Equal(t, "123456.123456", s.TriggerID)
		assert.Equal(t, "test-view-callback-id", s.View.CallbackID)
	})

	t.Run("Invalid view", func(t *testing.T) {
		requestData := url.Values{}
		requestData.Set("view", "invalid-view")
		err := s.ParseSlashCommand(requestData)
		assert.NotNil(t, err)
	})
}
