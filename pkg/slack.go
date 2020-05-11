package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func fetchChannelID(client *http.Client, token string, channel string) (string, error) {
	channelsInfoURL := fmt.Sprintf("https://slack.com/api/conversations.list?token=%s", token)
	cinfo, err := client.Get(channelsInfoURL)
	if err != nil {
		return "", err
	}
	defer cinfo.Body.Close()

	data, err := ioutil.ReadAll(cinfo.Body)
	if err != nil {
		return "", err
	}

	channelInfo := struct {
		OK       bool `json:"ok"`
		Channels []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"channels"`
	}{}

	if err := json.Unmarshal(data, &channelInfo); err != nil {
		return "", err
	}

	if !channelInfo.OK {
		return "", fmt.Errorf("failed to fetch channel info")
	}

	var channelID string
	for _, ch := range channelInfo.Channels {
		if ch.Name == channel {
			channelID = ch.ID
		}
	}

	if channelID == "" {
		return "", fmt.Errorf("cannot find channel with name %s", channel)
	}

	return channelID, nil
}

// SendMessage sends a message to specific channel on Slack.
func SendMessage(channel, message, token string) error {
	client := slackConf.Client(context.Background(), &oauth2.Token{AccessToken: token})

	channelID, err := fetchChannelID(client, token, channel)
	if err != nil {
		return err
	}

	msgBody := struct {
		Channel string `json:"channel"`
		Text    string `json:"text"`
		AsUser  bool   `json:"as_user"`
	}{
		Channel: channelID,
		Text:    message,
		AsUser:  true,
	}

	messageBody, err := json.Marshal(msgBody)
	if err != nil {
		return err
	}

	buf := bytes.NewReader(messageBody)
	sinfo, err := client.Post("https://slack.com/api/chat.postMessage", "application/json", buf)
	if err != nil {
		return err
	}
	defer sinfo.Body.Close()

	data, err := ioutil.ReadAll(sinfo.Body)
	if err != nil {
		return err
	}

	sendInfo := struct {
		OK bool `json:"ok"`
	}{}

	if err := json.Unmarshal(data, &sendInfo); err != nil {
		return err
	}

	if !sendInfo.OK {
		return fmt.Errorf("failed to send message")
	}

	return nil
}

func fetchUsersMap(client *http.Client, token string) (map[string]string, error) {
	users := map[string]string{
		"": "slack",
	}

	usersInfoURL := fmt.Sprintf("https://slack.com/api/users.list?token=%s", token)
	uinfo, err := client.Get(usersInfoURL)
	if err != nil {
		return nil, err
	}
	defer uinfo.Body.Close()

	data, err := ioutil.ReadAll(uinfo.Body)
	if err != nil {
		return nil, err
	}

	usersInfo := struct {
		OK    bool `json:"ok"`
		Users []struct {
			ID   string `json:"ID"`
			Name string `json:"name"`
		} `json:"members"`
	}{}

	if err := json.Unmarshal(data, &usersInfo); err != nil {
		return nil, err
	}

	if !usersInfo.OK {
		return nil, fmt.Errorf("failed to fetch users info")
	}

	for _, u := range usersInfo.Users {
		users[u.ID] = u.Name
	}

	return users, nil
}

// SlackMessage is an instance of mwssage from slack.
type SlackMessage struct {
	User      string `json:"user"`
	Text      string `json:"text"`
	TimeStamp string `json:"ts"`
}

// FetchMessages gets the unread messages from slack.
func FetchMessages(channel string, limit uint, token string) ([]*SlackMessage, error) {
	client := slackConf.Client(context.Background(), &oauth2.Token{AccessToken: token})

	channelID, err := fetchChannelID(client, token, channel)
	if err != nil {
		return nil, err
	}

	users, err := fetchUsersMap(client, token)
	if err != nil {
		return nil, err
	}

	messageLimit := uint(20)
	if limit != 0 {
		messageLimit = limit
	}

	messageInfoURL := fmt.Sprintf("https://slack.com/api/conversations.history?token=%s&channel=%s&limit=%d", token, channelID, messageLimit)
	minfo, err := client.Get(messageInfoURL)
	if err != nil {
		return nil, err
	}
	defer minfo.Body.Close()

	data, err := ioutil.ReadAll(minfo.Body)
	if err != nil {
		return nil, err
	}

	messagesInfo := struct {
		OK       bool            `json:"ok"`
		Messages []*SlackMessage `json:"messages"`
	}{}

	if err := json.Unmarshal(data, &messagesInfo); err != nil {
		return nil, err
	}

	if !messagesInfo.OK {
		return nil, fmt.Errorf("failed to fetch messages")
	}

	for _, m := range messagesInfo.Messages {
		m.User = users[m.User]
	}

	return messagesInfo.Messages, nil
}
