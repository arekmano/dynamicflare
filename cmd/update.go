package cmd

import (
	"github.com/arekmano/dynamicflare/service"

	"github.com/spf13/cobra"
)

var isDryRun bool

func init() {
	updateCmd.Flags().BoolVarP(&isDryRun, "dryrun", "d", false, "if set to true, will not update the cloudflare entry")

}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "used to update cloudflare dynamically",
	Long:  "used to update cloudflare dynamically",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := setup(verbose, configFile)
		return service.
			New(config).
			Run(isDryRun, config.Records)
	},
}
