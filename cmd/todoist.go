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

var todoistAddCmd = &cobra.Command{
	Use:   "add [task content]",
	Short: "Add a task to your todoist",

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

		task := args[0]

		if err := pkg.AddTask(conf.TodoistToken, task); err != nil {
			logrus.WithError(err).Fatalln("Cannot add task")
		}

		logrus.Infoln("Task added successfully")
	},
}

var todoistTodayCmd = &cobra.Command{
	Use:   "today",
	Short: "Fetch today's tasks from todoist",

	Run: func(cmd *cobra.Command, args []string) {
		conf, err := pkg.LoadConfig()
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot load config")
		}

		tasks, err := pkg.GetTasks(conf.TodoistToken)
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot fetch tasks")
		}

		if len(tasks) == 0 {
			logrus.Warnln("No tasks for today")
			return
		}

		for _, t := range tasks {
			logrus.Infoln(t)
		}
	},
}
