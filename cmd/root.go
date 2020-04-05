package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dynamicflare",
	Short: "used to interact with cloudflare dynamically",
	Long:  "used to interact with cloudflare dynamically",
}

var configFile string
var verbose bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "the configuration file to use")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "if set to true, set the debug level to debug")
	rootCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(domainsCmd)
	rootCmd.AddCommand(recordsCmd)
}

// Execute executes the root cmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validate() error {
	if configFile == "" {
		return errors.New("Must specify the config file")
	}
	stat, err := os.Stat(configFile)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return errors.New("Specified a directory. Must specify a file")
	}
	return nil
}
