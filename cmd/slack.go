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
