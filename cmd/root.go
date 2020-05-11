package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shubhamgupta2956/mind/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mind",
	Short: "Handy tool to carry out your day to day work",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := pkg.CreateDirIfNotExist(pkg.ConfigDirPath); err != nil {
			logrus.WithError(err).Fatalln("Cannot create config directory")
		}

		confDefaults, err := defaultConfContent()
		if err != nil {
			confDefaults = `{"output_type": "plain"}`
		}
		if err := pkg.CreateFileIfNotExist(pkg.ConfigFilePath, confDefaults); err != nil {
			logrus.WithError(err).Fatalln("Cannot create config file")
		}

		todayDate := time.Now().Format(pkg.DateLayout)
		historyPath, err := pkg.ConfigHistoryPath(todayDate)
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot get history filepath")
		}

		if err := pkg.CreateFileIfNotExist(historyPath, "time,cmd\n"); err != nil {
			logrus.WithError(err).Fatalln("Cannot create today's history file")
		}

		if err := updateOutputType(); err != nil {
			logrus.WithError(err).Fatalln("Cannot update output type")
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logrus.WithError(err).Fatalln("Cannot execute cmd")
		}
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		message := getFinalCmd(cmd, args)

		todayDate := time.Now().Format(pkg.DateLayout)
		historyPath, err := pkg.ConfigHistoryPath(todayDate)
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot get history filepath")
		}

		if err := logToCSV(message, historyPath); err != nil {
			logrus.WithError(err).Fatalln("Cannot log command to history")
		}
	},
}

func getFinalCmd(cmd *cobra.Command, args []string) string {
	if cmd.HasParent() {
		prevCmd := getFinalCmd(cmd.Parent(), nil)
		joinArgs := []string{prevCmd}
		joinArgs = append(joinArgs, cmd.CalledAs())
		joinArgs = append(joinArgs, args...)
		return strings.Join(joinArgs, " ")
	}

	return "mind"
}

func logToCSV(msg, file string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line := fmt.Sprintf("%s,%s\n", time.Now().Format("15:04:05"), msg)
	if _, err := f.WriteString(line); err != nil {
		return err
	}

	return nil
}

func defaultConfContent() (string, error) {
	conf := pkg.Config{
		OutputType: "plain",
	}

	content, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return "", err
	}

	return string(content), nil
}
