package helpers

import (
	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/spf13/cobra"
)

// MakeNodeCompletion is a wrapper for `cobra.Command.RegisterFlagCompletionFunc`
// to reuse in commands that want to have `--node` flag completion
func MakeNodeCompletion(cmd *cobra.Command) error {
	return cmd.RegisterFlagCompletionFunc("node", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		var nodes []string
		thCfg, _ := cmd.Flags().GetString("config-file")
		thEnvFiles, _ := cmd.Flags().GetStringSlice("env-file")

		cfg, err := config.LoadAndValidateFromFile(thCfg, thEnvFiles, false)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		for i := range cfg.Nodes {
			nodes = append(nodes, cfg.Nodes[i].Hostname)
			nodes = append(nodes, cfg.Nodes[i].IPAddress)
		}
		return nodes, cobra.ShellCompDirectiveNoFileComp
	})
}
