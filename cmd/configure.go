package cmd

import (
	"github.com/shubhamgupta2956/mind/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure [output configuration]",
	Short: "Configure type of output to screen",
	Long:  "Log output configuration, can be either 'json' or 'plain'.",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logrus.Fatalln("Invalid number of arguments. Requires one.")
		}

		switch args[0] {
		case "plain", "json":
		default:
			logrus.WithField("output_type", args[0]).Fatalf("Output type can only be 'plain' or 'json'")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := pkg.LoadConfig()
		if err != nil {
			logrus.WithError(err).Fatalln("Cannot load config")
		}

		conf.OutputType = args[0]
		if err := pkg.UpdateConfig(conf); err != nil {
			logrus.WithError(err).Fatalln("Cannot update config")
		}

		if err := updateOutputType(); err != nil {
			logrus.WithError(err).Fatalln("Cannot update output type")
		}

		logrus.WithField("output_type", args[0]).Infof("Output type changed successfully")
	},
}

func updateOutputType() error {
	conf, err := pkg.LoadConfig()
	if err != nil {
		return err
	}

	if conf.OutputType == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	return nil
}
