package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// DateLayout is the layout in which to parse the date.
const DateLayout = "2006-01-02"

var (
	// ConfigDirPath is the path to mind config directory ($HOME/.mind)
	ConfigDirPath string

	// ConfigFilePath is the path to mind configuration file.
	ConfigFilePath string
)

// Config represents the configuration file.
type Config struct {
	OutputType   string `json:"output_type,omitempty"`
	SlackToken   string `json:"slack_token,omitempty"`
	TodoistToken string `json:"todoist_token,omitempty"`
}

func init() {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		logrus.Fatalln("No $HOME found")
	}

	ConfigDirPath = filepath.Join(home, ".mind")
	ConfigFilePath = filepath.Join(ConfigDirPath, "config.json")
}

// ConfigHistoryPath returns path to config file corresponding to the date.
// Takes date in the format "YYYY-MM-DD". Returns error if the date is invalid.
func ConfigHistoryPath(date string) (string, error) {
	t, err := time.Parse(DateLayout, date)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.%s", t.Format(DateLayout), "csv")
	path := filepath.Join(ConfigDirPath, filename)
	return path, nil
}

// LoadConfig loads the config data.
func LoadConfig() (Config, error) {
	file, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		return Config{}, err
	}

	var conf Config
	if err := json.Unmarshal(file, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}

// UpdateConfig updates the config with the new config.
func UpdateConfig(update Config) error {
	file, err := json.MarshalIndent(update, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(ConfigFilePath, file, 0644)
}
