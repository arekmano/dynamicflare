package cmd

import (
	"github.com/arekmano/dynamicflare/service"

	"github.com/spf13/cobra"
)

var recordsCmd = &cobra.Command{
	Use:   "records",
	Short: "used to print cloudflare records dynamically",
	Long:  "used to print cloudflare records dynamically",
	RunE: func(cmd *cobra.Command, args []string) error {
		return service.
			New(setup(verbose, configFile)).
			ListDomainRecords()
	},
}
