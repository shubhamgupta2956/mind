package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mind",
	Short: "Handy tool to carry out your day to day work",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logrus.WithError(err).Fatalln("Cannot execute cmd")
		}
	},
}
