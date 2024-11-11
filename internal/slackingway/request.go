package slackingway

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

// Body of request data from Slack
type SlackRequestBody struct {
	Timestamp   string     `json:"timestamp"`
	Type        string     `json:"type"`
	Challenge   string     `json:"challenge"`
	Token       string     `json:"token"`
	Command     string     `json:"command"`
	Text        string     `json:"text"`
	ResponseURL string     `json:"response_url"`
	UserID      string     `json:"user_id"`
	ChannelID   string     `json:"channel_id"`
	TeamID      string     `json:"team_id"`
	CallbackID  string     `json:"callback_id"`
	TriggerID   string     `json:"trigger_id"`
	View        slack.View `json:"view"`
	Event       SlackEvent `json:"event"`
}

// ParseTimestamp parses the timestamp from Slack
func (s *SlackRequestBody) ParseTimestamp(timestamp string) error {
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Printf("Error parsing timestamp: %v", err)
		return err
	}

	s.Timestamp = time.Unix(timestampInt, 0).UTC().Format(time.RFC822)

	return nil
}

// ParsePayload parses the payload from Slack
func (s *SlackRequestBody) ParsePayload(payload string) error {
	// Parse the payload
	if err := json.Unmarshal([]byte(payload), s); err != nil {
		log.Printf("Error parsing payload: %v", err)
		return err
	}
	s.Command = s.View.CallbackID
	s.CallbackID = s.View.CallbackID

	// Parse the payload as a map
	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
		log.Printf("Error parsing payload: %v", err)
		return err
	}

	// Get User ID
	user, ok := payloadMap["user"].(map[string]interface{})
	if !ok {
		log.Println("User field not found")
		return errors.New("User field not found")
	}
	userID, ok := user["id"].(string)
	if !ok {
		log.Println("User ID field not found")
		return errors.New("User ID field not found")
	}
	s.UserID = userID

	return nil
}

// ParseSlashCommand parses the slash command request from Slack
func (s *SlackRequestBody) ParseSlashCommand(requestData url.Values) error {
	s.Type = "slash_command"
	s.Token = requestData.Get("token")
	s.Command = requestData.Get("command")
	s.Text = requestData.Get("text")
	s.ResponseURL = requestData.Get("response_url")
	s.UserID = requestData.Get("user_id")
	s.ChannelID = requestData.Get("channel_id")
	s.TeamID = requestData.Get("team_id")
	s.CallbackID = requestData.Get("callback_id")
	s.TriggerID = requestData.Get("trigger_id")

	// Parse the view
	if requestData.Get("view") != "" {
		var view slack.View
		if err := json.Unmarshal([]byte(requestData.Get("view")), &view); err != nil {
			log.Printf("Error unmarshaling view: %v", err)
			return err
		}
		s.View = view
	}

	return nil
}
