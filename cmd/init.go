package cmd

import "github.com/sirupsen/logrus"

func init() {
	rootCmd.AddCommand(
		configureCmd,
	)
}

// Execute starts the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatalln("cannot execute cmd")
	}
}
