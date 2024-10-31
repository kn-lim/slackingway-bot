package slackingway

import "github.com/slack-go/slack"

func Echo(s *SlackingwayWrapper) error {
	// Generate a new modal
	modal := GenerateEchoModal()

	// Open a new modal
	_, err := s.APIClient.OpenView(s.SlackRequestBody.TriggerID, modal)
	if err != nil {
		return err
	}

	return nil
}

func GenerateEchoModal() slack.ModalViewRequest {
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

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			inputBlock,
		},
	}

	return slack.ModalViewRequest{
		Type:       slack.ViewType("modal"),
		Title:      titleText,
		Close:      closeText,
		Submit:     submitText,
		Blocks:     blocks,
		CallbackID: "echo",
	}
}

func UpdateEchoModal() slack.ModalViewRequest {
	// Create a new Slack ModalViewRequest
	titleText := slack.NewTextBlockObject("plain_text", "Echo", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)
	headerText := slack.NewTextBlockObject("mrkdwn", "Modal updated!", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
		},
	}

	return slack.ModalViewRequest{
		Type:       slack.ViewType("modal"),
		Title:      titleText,
		Close:      closeText,
		Submit:     submitText,
		Blocks:     blocks,
		CallbackID: "echo",
	}
}
