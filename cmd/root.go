package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var rootLongHelp = strings.TrimSpace(`
talhelper is a tool to help you create a Talos cluster.

Currently there is only one usage, which is creating a cluster config YAML file.

Workflow:
	taloshelper genconfig --config-file config.yaml --out-dir ./clusterconfig --vars-file vars.yaml
`)

var rootCmd = &cobra.Command{
	Use: "talhelper",
	Short: "A tool to help with creating Talos cluster",
	Long: rootLongHelp,
	SilenceUsage: true,
	SilenceErrors: true,
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}
