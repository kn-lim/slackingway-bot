package slackingway_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

func TestValidateRequest(t *testing.T) {
	testSlackSigningSecret := "test-slack-signing-secret"

	// Create a mock request
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	body := "test-body"
	basestring := fmt.Sprintf("%s:%s:%s", slackingway.SlackVersion, timestamp, body)
	h := hmac.New(sha256.New, []byte(testSlackSigningSecret))
	h.Write([]byte(basestring))
	signature := slackingway.SlackVersion + "=" + hex.EncodeToString(h.Sum(nil))

	request := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"X-Slack-Request-Timestamp": timestamp,
			"X-Slack-Signature":         signature,
		},
		Body: body,
	}

	assert.Nil(t, slackingway.ValidateRequest(request, testSlackSigningSecret))
}

func TestValidateRole(t *testing.T) {
	t.Run("Empty role users", func(t *testing.T) {
		testRole := ""

		assert.False(t, slackingway.ValidateRole(testRole, "test-user1"))
	})

	t.Run("Valid role users", func(t *testing.T) {
		testRole := "test-user1,test-user2"
		assert.True(t, slackingway.ValidateRole(testRole, "test-user1"))
		assert.True(t, slackingway.ValidateRole(testRole, "test-user2"))
		assert.False(t, slackingway.ValidateRole(testRole, "test-user3"))
	})
}
