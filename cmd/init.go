package cmd

import "github.com/sirupsen/logrus"

func init() {
	slackCmd.AddCommand(
		slackAuthCmd,
		slackSendCmd,
	)

	todoistCmd.AddCommand(
		todoistAuthCmd,
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
