package slackingway

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
