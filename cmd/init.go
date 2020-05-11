package cmd

import (
	"github.com/sirupsen/logrus"
)

var (
	messageLimit uint
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

// Execute starts the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatalln("cannot execute cmd")
	}
}
