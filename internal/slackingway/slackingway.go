package slackingway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

// Body of request data from Slack
type SlackRequestBody struct {
	Timestamp   string `json:"timestamp"`
	Type        string `json:"type"`
	Challenge   string `json:"challenge"`
	Token       string `json:"token"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
	UserID      string `json:"user_id"`
	ChannelID   string `json:"channel_id"`
	TeamID      string `json:"team_id"`
}

type SlackAPIClient interface {
	GetUserInfo(userID string) (*slack.User, error)
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

type Slackingway interface {
	// Slack Specific
	NewResponse(message slack.Msg) (*http.Request, error)
	SendResponse(request *http.Request) error
	WriteToHistory() error
}

type SlackingwayWrapper struct {
	APIClient        SlackAPIClient
	HTTPClient       *http.Client
	SlackRequestBody *SlackRequestBody
}

// NewSlackingway creates a new SlackingwayWrapper
func NewSlackingway(s *SlackRequestBody) *SlackingwayWrapper {
	return &SlackingwayWrapper{
		APIClient:        slack.New(os.Getenv("SLACK_OAUTH_TOKEN")),
		HTTPClient:       &http.Client{},
		SlackRequestBody: s,
	}
}

// NewResponse creates a new HTTP request for a Slack response
func (s *SlackingwayWrapper) NewResponse(message slack.Msg) (*http.Request, error) {
	// Convert the response to JSON
	responseBody, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return nil, err
	}

	// Create the request
	request, err := http.NewRequest("POST", s.SlackRequestBody.ResponseURL, bytes.NewBuffer(responseBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

// SendResponse sends a response to Slack
func (s *SlackingwayWrapper) SendResponse(request *http.Request) error {
	// Create HTTP client and send request
	response, err := s.HTTPClient.Do(request)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return err
	}
	defer response.Body.Close()

	// Log the response status and body
	log.Printf("Response status: %v", response.Status)
	responseBodyBytes, _ := io.ReadAll(response.Body)
	log.Printf("Response body: %v", string(responseBodyBytes))

	// Check for non-OK status
	if response.StatusCode != http.StatusOK {
		log.Printf("Non-OK HTTP status: %v", response.StatusCode)
		return err
	}

	return nil
}

// WriteToHistory writes a message to Slackingway's History channel
func (s *SlackingwayWrapper) WriteToHistory() error {
	// Check if APIClient is nil
	if s.APIClient == nil {
		log.Printf("APIClient is nil")
		return fmt.Errorf("APIClient is nil")
	}

	// Get user information
	user, err := s.APIClient.GetUserInfo(s.SlackRequestBody.UserID)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return err
	}

	// Post a message to the History channel
	full_command := s.SlackRequestBody.Command
	if s.SlackRequestBody.Text != "" {
		full_command += " " + s.SlackRequestBody.Text
	}
	msg := fmt.Sprintf("%s executed command `%s` on %s", user.RealName, full_command, s.SlackRequestBody.Timestamp)
	_, _, err = s.APIClient.PostMessage(
		os.Getenv("SLACK_HISTORY_CHANNEL_ID"),
		slack.MsgOptionText(msg, false),
	)
	if err != nil {
		log.Printf("Error posting message: %v", err)
		return err
	}

	return nil
}
