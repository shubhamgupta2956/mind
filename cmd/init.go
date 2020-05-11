package cmd

import "github.com/sirupsen/logrus"

func init() {
	slackCmd.AddCommand(
		slackAuthCmd,
	)

	rootCmd.AddCommand(
		configureCmd,
		slackCmd,
	)
}

// Execute starts the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatalln("cannot execute cmd")
	}
}
