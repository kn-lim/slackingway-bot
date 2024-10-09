package slackingway

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

// NewResponse creates a new HTTP request for a Slack response
func NewResponse(responseURL string, message slack.Msg) (*http.Request, error) {
	// Convert the response to JSON
	responseBody, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return nil, err
	}

	// Create the request
	request, err := http.NewRequest("POST", responseURL, bytes.NewBuffer(responseBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

// SendResponse sends a response to Slack
func SendResponse(request *http.Request) error {
	// Create HTTP client and send request
	client := &http.Client{}
	response, err := client.Do(request)
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
