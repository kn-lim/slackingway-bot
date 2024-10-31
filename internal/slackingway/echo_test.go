package slackingway_test

import (
	"testing"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEcho(t *testing.T) {
	// Create a new mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	MockSlackAPIClient := NewMockSlackAPIClient(ctrl)
	MockSlackAPIClient.EXPECT().OpenView(gomock.Any(), gomock.Any()).Return(&slack.ViewResponse{}, nil)

	s := &slackingway.SlackingwayWrapper{
		APIClient:        MockSlackAPIClient,
		SlackRequestBody: &slackingway.SlackRequestBody{TriggerID: "definitely_a_valid_trigger_id"},
	}
	err := slackingway.Echo(s)

	// Run test
	assert.Nil(t, err)
}

func TestGenerateEchoModal(t *testing.T) {
	// Create expected Slack ModalViewRequest
	expectedTitle := slack.NewTextBlockObject("plain_text", "Echo", false, false)
	expectedClose := slack.NewTextBlockObject("plain_text", "Close", false, false)
	expectedSubmit := slack.NewTextBlockObject("plain_text", "Submit", false, false)
	expectedHeaderText := slack.NewTextBlockObject("mrkdwn", "Enter the text you want to echo", false, false)
	expectedInputText := slack.NewTextBlockObject("plain_text", "Text", false, false)
	expectedInputHint := slack.NewTextBlockObject("plain_text", "Enter the text you want to echo", false, false)
	expectedInputPlaceholder := slack.NewTextBlockObject("plain_text", "Enter text here", false, false)
	expectedInputElement := slack.NewPlainTextInputBlockElement(expectedInputPlaceholder, "text")

	expectedBlocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(expectedHeaderText, nil, nil),
			slack.NewInputBlock("text", expectedInputText, expectedInputHint, expectedInputElement),
		},
	}

	expected := slack.ModalViewRequest{
		Type:       slack.ViewType("modal"),
		Title:      expectedTitle,
		Close:      expectedClose,
		Submit:     expectedSubmit,
		Blocks:     expectedBlocks,
		CallbackID: "echo",
	}

	got := slackingway.GenerateEchoModal()

	// Run test
	assert.Equal(t, expected, got)
}

func TestUpdateEchoModal(t *testing.T) {
	// Create expected Slack ModalViewRequest
	expectedTitle := slack.NewTextBlockObject("plain_text", "Echo", false, false)
	expectedClose := slack.NewTextBlockObject("plain_text", "Close", false, false)
	expectedSubmit := slack.NewTextBlockObject("plain_text", "Submit", false, false)
	expectedHeaderText := slack.NewTextBlockObject("mrkdwn", "Modal updated!", false, false)
	expectedHeaderSection := slack.NewSectionBlock(expectedHeaderText, nil, nil)

	expectedBlocks := slack.Blocks{
		BlockSet: []slack.Block{
			expectedHeaderSection,
		},
	}

	expected := slack.ModalViewRequest{
		Type:       slack.ViewType("modal"),
		Title:      expectedTitle,
		Close:      expectedClose,
		Submit:     expectedSubmit,
		Blocks:     expectedBlocks,
		CallbackID: "echo",
	}

	got := slackingway.UpdateEchoModal()

	// Run test
	assert.Equal(t, expected, got)
}
