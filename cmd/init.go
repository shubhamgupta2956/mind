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
		todoistTodayCmd,
	)

	sessionCmd.AddCommand(
		sessionSummaryCmd,
	)

	rootCmd.AddCommand(
		configureCmd,
		slackCmd,
		todoistCmd,
		sessionCmd,
	)
}

// Execute starts the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatalln("cannot execute cmd")
	}
}
