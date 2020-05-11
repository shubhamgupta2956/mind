package cmd

import (
	"github.com/shubhamgupta2956/mind/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Authorize, fetch and send messages through slack",

	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logrus.WithError(err).Fatalln("Cannot run slack command")
		}
	},
}

var slackAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authorize slack workspace for mind",

	Run: func(cmd *cobra.Command, args []string) {
		conf, err := pkg.LoadConfig()
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot load config")
		}

		s := pkg.Slack{}
		pkg.RunAuthCmd(&s, &conf.SlackToken)

		if err := pkg.UpdateConfig(conf); err != nil {
			logrus.WithError(err).Fatalln("Cannot update config")
		}
	},
}

var slackSendCmd = &cobra.Command{
	Use:   "send [channel name] [message]",
	Short: "Send message to slack workspace channel",

	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logrus.Fatalln("Invalid number of arguments. Requires two.")
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		conf, err := pkg.LoadConfig()
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot load config")
		}

		channel := args[0]
		message := args[1]

		if err := pkg.SendMessage(channel, message, conf.SlackToken); err != nil {
			logrus.WithError(err).Fatalln("Cannot fetch messages")
		}

		logrus.Infoln("Message sent successfully")
	},
}

var slackUnreadCmd = &cobra.Command{
	Use:   "unreads [channel name]",
	Short: "Get unread messages from slack workspace",

	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logrus.Fatalln("Invalid number of arguments. Requires one.")
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		conf, err := pkg.LoadConfig()
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot load config")
		}

		channel := args[0]

		messages, err := pkg.FetchMessages(channel, messageLimit, conf.SlackToken)
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot fetch messages")
		}

		if len(messages) == 0 {
			logrus.Warnln("No messages found in the cannel")
			return
		}

		for _, message := range messages {
			logrus.WithField("user", message.User).WithField(
				"timestamp", message.TimeStamp).Infoln(message.Text)
		}
	},
}
