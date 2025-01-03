package slackingway

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

func Menu(s *SlackingwayWrapper) error {
	// Create a new modal
	modal := CreateMenuModal()

	// Open a new modal
	_, err := s.APIClient.OpenView(s.SlackRequestBody.TriggerID, modal)
	if err != nil {
		return err
	}

	return nil
}

func ReceivedMenu(s *SlackingwayWrapper) ([]slack.Block, error) {
	// Get State
	state := s.SlackRequestBody.View.State.Values

	// Get the selected options
	option1 := state["action_option1"]["menu_option1"].SelectedOption.Value
	if option1 == "" {
		option1 = "N/A"
	}
	option2 := state["action_option2"]["menu_option2"].SelectedOption.Value
	if option2 == "" {
		option2 = "N/A"
	}
	option3 := state["action_option3"]["menu_option3"].SelectedOption.Value
	if option3 == "" {
		option3 = "N/A"
	}

	// Get user information
	user, err := s.APIClient.GetUserInfo(s.SlackRequestBody.UserID)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return nil, err
	}

	// Create message blocks
	titleText := slack.NewTextBlockObject("plain_text", fmt.Sprintf("Received Menu Selections from %s", user.RealName), false, false)
	titleSection := slack.NewSectionBlock(titleText, nil, nil)

	option1Text := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Option 1:* `%s`", option1), false, false)
	option1Section := slack.NewSectionBlock(option1Text, nil, nil)

	option2Text := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Option 2:* `%s`", option2), false, false)
	option2Section := slack.NewSectionBlock(option2Text, nil, nil)

	option3Text := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Option 3:* `%s`", option3), false, false)
	option3Section := slack.NewSectionBlock(option3Text, nil, nil)

	return []slack.Block{
		titleSection,
		slack.NewDividerBlock(),
		option1Section,
		option2Section,
		option3Section,
	}, nil
}

func CreateMenuModal() slack.ModalViewRequest {
	// Create a new Slack ModalViewRequest
	titleText := slack.NewTextBlockObject("plain_text", "Menu", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Select the options", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Option 1
	option1Text := slack.NewTextBlockObject("plain_text", "Option 1", false, false)
	option1Section := slack.NewSectionBlock(option1Text, nil, nil)
	option1Options := []*slack.OptionBlockObject{
		{
			Text:  slack.NewTextBlockObject("plain_text", "Choice 1", false, false),
			Value: "choice1",
		},
		{
			Text:  slack.NewTextBlockObject("plain_text", "Choice 2", false, false),
			Value: "choice2",
		},
		{
			Text:  slack.NewTextBlockObject("plain_text", "Choice 3", false, false),
			Value: "choice3",
		},
	}
	option1ChoiceText := slack.NewTextBlockObject("plain_text", "Select a choice", false, false)
	option1Select := slack.NewOptionsSelectBlockElement("static_select", option1ChoiceText, "menu_option1", option1Options...)

	// Option 2
	option2Text := slack.NewTextBlockObject("plain_text", "Option 2", false, false)
	option2Section := slack.NewSectionBlock(option2Text, nil, nil)
	option2Options := []*slack.OptionBlockObject{
		{
			Text:  slack.NewTextBlockObject("plain_text", "Choice 1", false, false),
			Value: "choice1",
		},
		{
			Text:  slack.NewTextBlockObject("plain_text", "Choice 2", false, false),
			Value: "choice2",
		},
		{
			Text:  slack.NewTextBlockObject("plain_text", "Choice 3", false, false),
			Value: "choice3",
		},
	}
	option2ChoiceText := slack.NewTextBlockObject("plain_text", "Select a choice", false, false)
	option2Select := slack.NewOptionsSelectBlockElement("static_select", option2ChoiceText, "menu_option2", option2Options...)

	// Option 3
	option3Text := slack.NewTextBlockObject("plain_text", "Option 3", false, false)
	option3Section := slack.NewSectionBlock(option3Text, nil, nil)
	option3Options := []*slack.OptionBlockObject{
		{
			Text:  slack.NewTextBlockObject("plain_text", "Choice 1", false, false),
			Value: "choice1",
		},
		{
			Text:  slack.NewTextBlockObject("plain_text", "Choice 2", false, false),
			Value: "choice2",
		},
		{
			Text:  slack.NewTextBlockObject("plain_text", "Choice 3", false, false),
			Value: "choice3",
		},
	}
	option3ChoiceText := slack.NewTextBlockObject("plain_text", "Select a choice", false, false)
	option3Select := slack.NewOptionsSelectBlockElement("static_select", option3ChoiceText, "menu_option3", option3Options...)

	// Combine all the blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			slack.NewDividerBlock(),
			option1Section,
			slack.NewActionBlock("action_option1", option1Select),
			slack.NewDividerBlock(),
			option2Section,
			slack.NewActionBlock("action_option2", option2Select),
			slack.NewDividerBlock(),
			option3Section,
			slack.NewActionBlock("action_option3", option3Select),
		},
	}

	return slack.ModalViewRequest{
		Type:       slack.VTModal,
		Title:      titleText,
		Close:      closeText,
		Submit:     submitText,
		Blocks:     blocks,
		CallbackID: "/menu",
	}
}
