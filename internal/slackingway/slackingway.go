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

type SlackAPIClient interface {
	GetUserInfo(userID string) (*slack.User, error)
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
	OpenView(triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error)
	UpdateView(view slack.ModalViewRequest, externalID string, hash string, viewID string) (*slack.ViewResponse, error)
	PublishView(userID string, view slack.HomeTabViewRequest, hash string) (*slack.ViewResponse, error)
}

type Slackingway interface {
	NewResponse(message slack.Msg) (*http.Request, error)
	SendResponse(request *http.Request) error
	SendTextMessage(message, channelID string) error
	SendBlockMessage(blocks []slack.Block, channelID string) error
	WriteToHistory() error
}

type SlackingwayWrapper struct {
	Debug            bool
	APIClient        SlackAPIClient
	HTTPClient       *http.Client
	SlackRequestBody *SlackRequestBody
}

// NewSlackingway creates a new SlackingwayWrapper
func NewSlackingway(s *SlackRequestBody) *SlackingwayWrapper {
	return &SlackingwayWrapper{
		Debug:            os.Getenv("DEBUG") == "true",
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
	if s.Debug {
		log.Printf("Response status: %v", response.Status)
		responseBodyBytes, _ := io.ReadAll(response.Body)
		log.Printf("Response body: %v", string(responseBodyBytes))
	}

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

// SendTextMessage sends a message to a Slack channel
func (s *SlackingwayWrapper) SendTextMessage(message, channelID string) error {
	// Check if APIClient is nil
	if s.APIClient == nil {
		log.Printf("APIClient is nil")
		return fmt.Errorf("APIClient is nil")
	}

	// Send the message to the channel
	_, _, err := s.APIClient.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	return nil
}

// SendBlockMessage sends a message with blocks to a Slack channel
func (s *SlackingwayWrapper) SendBlockMessage(blocks []slack.Block, channelID string) error {
	// Check if APIClient is nil
	if s.APIClient == nil {
		log.Printf("APIClient is nil")
		return fmt.Errorf("APIClient is nil")
	}

	// Send the message to the channel
	_, _, err := s.APIClient.PostMessage(channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	return nil
}
