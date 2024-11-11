package slackingway

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

func Echo(s *SlackingwayWrapper) error {
	// Create a new modal
	modal := CreateEchoModal()

	// Open a new modal
	_, err := s.APIClient.OpenView(s.SlackRequestBody.TriggerID, modal)
	if err != nil {
		return err
	}

	return nil
}

func ReceivedEcho(s *SlackingwayWrapper) (slack.Msg, error) {
	// Get State
	state := s.SlackRequestBody.View.State.Values

	// Get the text from the input block
	text := state["text"]["text"].Value

	// Get user information
	user, err := s.APIClient.GetUserInfo(s.SlackRequestBody.UserID)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return slack.Msg{}, err
	}

	var message slack.Msg
	if text == "" {
		message = slack.Msg{
			Text: fmt.Sprintf("Received an empty echo from %s", user.RealName),
		}
	} else {
		message = slack.Msg{
			Text: fmt.Sprintf("Received Echo from %s: `%v`", user.RealName, text),
		}
	}

	return message, nil
}

func CreateEchoModal() slack.ModalViewRequest {
	// Create a new Slack ModalViewRequest
	titleText := slack.NewTextBlockObject("plain_text", "Echo", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Enter the text you want to echo", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	inputText := slack.NewTextBlockObject("plain_text", "Text", false, false)
	inputHint := slack.NewTextBlockObject("plain_text", "Enter the text you want to echo", false, false)
	inputPlaceholder := slack.NewTextBlockObject("plain_text", "Enter text here", false, false)
	inputElement := slack.NewPlainTextInputBlockElement(inputPlaceholder, "text")
	inputBlock := slack.NewInputBlock("text", inputText, inputHint, inputElement)

	// Combine all the blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			inputBlock,
		},
	}

	return slack.ModalViewRequest{
		Type:       slack.VTModal,
		Title:      titleText,
		Close:      closeText,
		Submit:     submitText,
		Blocks:     blocks,
		CallbackID: "/echo",
	}
}
