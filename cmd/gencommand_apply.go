package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/pkg/generate"
	"github.com/budimanjojo/talhelper/pkg/parse"
)

var (
	gencommandApplyOutDir     string
	gencommandApplyCfgFile    string
	gencommandApplyEnvFile    []string

	gencommandApplyNode             string
	gencommandApplyExtraFlags           []string
)

var gencommandApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Generate talosctl apply-config commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := parse.ParseConfig(gencommandUpgradeCfgFile, gencommandUpgradeEnvFile)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		err = generate.GenerateApplyCommand(&cfg, gencommandApplyOutDir, gencommandApplyNode, gencommandApplyExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl apply command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandApplyCmd)
	gencommandApplyCmd.Flags().StringVarP(&gencommandApplyCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	gencommandApplyCmd.Flags().StringVarP(&gencommandApplyOutDir, "out-dir", "o", "./clusterconfig", "Directory that contains the generated config files to apply.")
	gencommandApplyCmd.Flags().StringSliceVar(&gencommandApplyEnvFile, "env-file", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	gencommandApplyCmd.Flags().StringSliceVar(&gencommandApplyExtraFlags, "extra-flags", []string{}, "List of additional flags that will be injected into the generated commands.")
	gencommandApplyCmd.Flags().StringVarP(&gencommandApplyNode, "node", "n", "", "A specific node to generate the command for. If not specified, will generate for all nodes.")
}
