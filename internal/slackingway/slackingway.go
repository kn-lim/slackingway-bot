package slackingway

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

// Body of request data from Slack
type SlackRequestBody struct {
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

type Slackingway struct {
	HTTPClient       *http.Client
	SlackRequestBody *SlackRequestBody
}

func NewSlackingway(s *SlackRequestBody) *Slackingway {
	return &Slackingway{
		HTTPClient: &http.Client{},
	}
}

// NewResponse creates a new HTTP request for a Slack response
func (s *Slackingway) NewResponse(message slack.Msg) (*http.Request, error) {
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
func (s *Slackingway) SendResponse(request *http.Request) error {
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
