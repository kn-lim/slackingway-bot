package slackingway_test

import (
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

func TestPing(t *testing.T) {
	expected := slack.Msg{Text: "Pong!"}
	actual, err := slackingway.Ping()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
