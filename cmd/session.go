package cmd

import (
	"encoding/csv"
	"os"

	"github.com/shubhamgupta2956/mind/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Session summary for mind",

	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logrus.WithError(err).Fatalln("Cannot run session command")
		}
	},
}

var sessionSummaryCmd = &cobra.Command{
	Use:   "summary [date as YYYY-MM-DD]",
	Short: "Session summary for mind",

	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logrus.Fatalln("Invalid number of arguments. Requires one.")
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		historyPath, err := pkg.ConfigHistoryPath(args[0])
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot get history filepath")
		}

		f, err := os.Open(historyPath)
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot open history file")
		}
		defer f.Close()

		reader := csv.NewReader(f)
		allLines, err := reader.ReadAll()
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot read history file")
		}

		if len(allLines) == 0 {
			logrus.Warnln("No history for today")
			return
		}

		for _, line := range allLines[1:] {
			logrus.WithField("time", line[0]).Infoln(line[1])
		}
	},
}
