package cmd

import (
	"dynamicflare/service"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "used to update cloudflare dynamically",
	Long:  "used to update cloudflare dynamically",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validate(); err != nil {
			return err
		}
		if isDryRun {
			logrus.SetLevel(logrus.DebugLevel)
		}
		config := loadConfig(configFile)
		logrus.
			Debugf("loaded configuration file: %+v", config)
		return service.
			New(config).
			Run(isDryRun, config.Records)
	},
}
