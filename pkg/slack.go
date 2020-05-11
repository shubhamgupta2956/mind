package pkg

import (
	"context"

	"github.com/shubhamgupta2956/mind"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/slack"
)

var slackConf = &oauth2.Config{
	ClientID:     mind.SlackClientID,
	ClientSecret: mind.SlackClientSecret,
	Endpoint:     slack.Endpoint,
	RedirectURL:  redirectURL,
	Scopes: []string{
		"channels:read",
		"channels:history",
		"chat:write:user",
		"users:read",
	},
}

// Slack is used to interact with the Slack API.
type Slack struct{}

// GetAuthURL gets the auth URL for slack.
func (s *Slack) GetAuthURL(state string) string {
	return slackConf.AuthCodeURL(state)
}

// GetToken get's the slack API token.
func (s *Slack) GetToken(code string) (string, error) {
	token, err := slackConf.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
