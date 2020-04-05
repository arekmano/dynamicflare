package cmd

import (
	"github.com/arekmano/dynamicflare/service"

	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "used to print cloudflare domains dynamically",
	Long:  "used to print cloudflare domains dynamically",
	RunE: func(cmd *cobra.Command, args []string) error {
		return service.
			New(setup(verbose, configFile)).
			ListDomains()
	},
}
