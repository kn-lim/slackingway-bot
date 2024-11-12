package slackingway

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

const (
	SlackVersion  = "v0"
	TimeAllowance = 30 * time.Second
)

// ValidateRequest verifies the request from Slack
func ValidateRequest(request events.APIGatewayProxyRequest, slackSigningSecret string) error {
	// Check timing
	timestampInt, err := strconv.ParseInt(request.Headers["X-Slack-Request-Timestamp"], 10, 64)
	if err != nil {
		log.Printf("Error parsing timestamp: %v", err)
		return err
	}
	if float64(time.Now().Unix())-float64(timestampInt) > TimeAllowance.Seconds() {
		log.Printf("Timestamp is too old")
		return errors.New("timestamp is too old")
	}

	// Check signature
	basestring := fmt.Sprintf("%s:%s:%s", SlackVersion, request.Headers["X-Slack-Request-Timestamp"], request.Body)
	h := hmac.New(sha256.New, []byte(slackSigningSecret))
	h.Write([]byte(basestring))
	signature := SlackVersion + "=" + hex.EncodeToString(h.Sum(nil))
	if request.Headers["X-Slack-Signature"] != signature {
		log.Printf("Invalid signature")
		return errors.New("invalid signature")
	}

	return nil
}

// ValidateRole checks if the user ID is in the role
func ValidateRole(role, userID string) bool {
	// Check if the user has the admin role
	adminRoleUsers := strings.Split(role, ",")
	for _, user := range adminRoleUsers {
		if user == userID {
			return true
		}
	}
	return false
}
