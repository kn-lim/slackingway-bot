package slackingway_test

import (
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
