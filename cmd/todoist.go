package cmd

import (
	"github.com/shubhamgupta2956/mind/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var todoistCmd = &cobra.Command{
	Use:   "todoist",
	Short: "Authorize, add and view tasks through todoist",

	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logrus.WithError(err).Fatalln("Cannot run todoist command")
		}
	},
}

var todoistAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authorize todoist workspace for mind",

	Run: func(cmd *cobra.Command, args []string) {
		conf, err := pkg.LoadConfig()
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot load config")
		}

		s := pkg.Todoist{}
		pkg.RunAuthCmd(&s, &conf.TodoistToken)

		if err := pkg.UpdateConfig(conf); err != nil {
			logrus.WithError(err).Fatalln("Cannot update config")
		}
	},
}
