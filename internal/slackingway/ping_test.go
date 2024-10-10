package slackingway_test

import (
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

// TestPing tests the Ping function
func TestPing(t *testing.T) {
	// Run test
	expected := slack.Msg{Text: "Pong!"}
	actual, err := slackingway.Ping()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
