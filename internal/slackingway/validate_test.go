package slackingway_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

func TestValidateRequest(t *testing.T) {
	// Set up environment variables
	os.Setenv("SLACK_SIGNING_SECRET", "test-slack-signing-secret")

	// Create a mock request
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	body := "test-body"
	basestring := fmt.Sprintf("%s:%s:%s", slackingway.SlackVersion, timestamp, body)
	h := hmac.New(sha256.New, []byte(os.Getenv("SLACK_SIGNING_SECRET")))
	h.Write([]byte(basestring))
	signature := slackingway.SlackVersion + "=" + hex.EncodeToString(h.Sum(nil))

	request := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"X-Slack-Request-Timestamp": timestamp,
			"X-Slack-Signature":         signature,
		},
		Body: body,
	}

	assert.Nil(t, slackingway.ValidateRequest(request))
}

func TestValidateAdminRole(t *testing.T) {
	t.Run("Empty admin role users", func(t *testing.T) {
		os.Setenv("ADMIN_ROLE_USERS", "")
		assert.False(t, slackingway.ValidateAdminRole("test-user1"))
	})

	t.Run("Valid admin role users", func(t *testing.T) {
		os.Setenv("ADMIN_ROLE_USERS", "test-user1,test-user2")
		assert.True(t, slackingway.ValidateAdminRole("test-user1"))
		assert.True(t, slackingway.ValidateAdminRole("test-user2"))
		assert.False(t, slackingway.ValidateAdminRole("test-user3"))
	})
}
