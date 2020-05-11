package pkg

import (
	"context"

	"github.com/shubhamgupta2956/mind"
	"golang.org/x/oauth2"
)

var todoistConf = &oauth2.Config{
	ClientID:     mind.TodoistClientID,
	ClientSecret: mind.TodoistClientSecret,
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://todoist.com/oauth/authorize",
		TokenURL: "https://todoist.com/oauth/access_token",
	},
	RedirectURL: redirectURL,
	Scopes:      []string{"data:read_write"},
}

// Todoist is used to interact with the todoist API.
type Todoist struct{}

// GetAuthURL gets the auth URL for todoist.
func (t *Todoist) GetAuthURL(state string) string {
	return todoistConf.AuthCodeURL(state)
}

// GetToken get's the todoist API token.
func (t *Todoist) GetToken(code string) (string, error) {
	token, err := todoistConf.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
