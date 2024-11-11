package slackingway

type SlackEvent struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	Tab     string `json:"tab"`
}
