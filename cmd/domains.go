package cmd

import (
	"dynamicflare/service"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "used to print cloudflare domains dynamically",
	Long:  "used to print cloudflare domains dynamically",
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
			ListDomains()
	},
}
