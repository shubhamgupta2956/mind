package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

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

func genTodoistCommands(uuid, typ string, args map[string]interface{}) (string, error) {
	commands := []struct {
		Type   string                 `json:"type"`
		TempID string                 `json:"temp_id"`
		Args   map[string]interface{} `json:"args"`
		UUID   string                 `json:"uuid"`
	}{{
		Type:   typ,
		TempID: "c7beb07f-b226-4eb1-bf63-30d782b07b1a", // Since we only do one command per request, this can be same
		UUID:   uuid,
		Args:   args,
	}}

	cmds, err := json.Marshal(commands)
	if err != nil {
		return "", err
	}

	return string(cmds), nil
}

func makeTodoistPostRequest(token, uuid, typ string, args map[string]interface{}) ([]byte, error) {
	client := todoistConf.Client(context.Background(), &oauth2.Token{AccessToken: token})

	val := url.Values{}

	if typ == "" || args == nil {
		val = url.Values{
			"token":          {token},
			"resource_types": {`["items"]`},
			"sync_token":     {"*"},
		}
	} else {
		commands, err := genTodoistCommands(uuid, typ, args)
		if err != nil {
			return nil, err
		}

		val = url.Values{
			"token":    {token},
			"commands": {commands},
		}
	}

	info, err := client.PostForm("https://api.todoist.com/sync/v8/sync", val)
	if err != nil {
		return nil, err
	}
	defer info.Body.Close()

	data, err := ioutil.ReadAll(info.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// AddTask adds a task to user's todoist.
func AddTask(token, task string) error {
	typ := "item_add"
	args := map[string]interface{}{
		"content": task,
		"due": map[string]interface{}{
			"date": time.Now().Format(DateLayout),
		},
	}

	uuid := randomToken()

	data, err := makeTodoistPostRequest(token, uuid, typ, args)
	if err != nil {
		return err
	}

	resp := struct {
		SyncStatus map[string]interface{} `json:"sync_status"`
	}{}

	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	if v, ok := resp.SyncStatus[uuid].(string); !ok || v != "ok" {
		return fmt.Errorf("error performing task")
	}

	return nil
}

// GetTasks gets today's tasks from todoist.
func GetTasks(token string) ([]string, error) {
	uuid := randomToken()

	data, err := makeTodoistPostRequest(token, uuid, "", nil)
	if err != nil {
		return nil, err
	}

	resp := struct {
		Items []struct {
			Content string      `json:"content"`
			Due     interface{} `json:"due"`
		} `json:"items"`
	}{}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	dateToday := time.Now().Format(DateLayout)
	tasks := []string{}

	for _, item := range resp.Items {
		dueDate, ok := item.Due.(map[string]interface{})
		if ok && strings.HasPrefix(dueDate["date"].(string), dateToday) {
			tasks = append(tasks, item.Content)
		}
	}

	return tasks, nil
}
