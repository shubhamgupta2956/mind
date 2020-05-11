package cmd

import "github.com/sirupsen/logrus"

var (
	messageLimit uint
)

func init() {
	slackUnreadCmd.Flags().UintVarP(&messageLimit, "limit", "l", 20, "Limit messages to this number")

	slackCmd.AddCommand(
		slackAuthCmd,
		slackSendCmd,
		slackUnreadCmd,
	)

	todoistCmd.AddCommand(
		todoistAuthCmd,
		todoistAddCmd,
	)

	rootCmd.AddCommand(
		configureCmd,
		slackCmd,
		todoistCmd,
	)
}

// Execute starts the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatalln("cannot execute cmd")
	}
}
