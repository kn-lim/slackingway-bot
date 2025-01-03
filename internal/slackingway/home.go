package slackingway

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

func HomeTab(s *SlackingwayWrapper) error {
	// Get user information
	user, err := s.APIClient.GetUserInfo(s.SlackRequestBody.Event.User)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return err
	}

	homeTabView := slack.HomeTabViewRequest{
		Type: slack.VTHomeTab,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Hello! :wave:\nWelcome to the Home Tab, %s", user.RealName), false, false),
					nil,
					nil,
				),
			},
		},
	}

	// Publish the view to the user's Home tab
	_, err = s.APIClient.PublishView(s.SlackRequestBody.Event.User, homeTabView, "")
	if err != nil {
		log.Printf("Error publishing Home tab view for %s: %v", user.RealName, err)
		return err
	}

	return nil
}
