package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/pkg/generate"
	"github.com/budimanjojo/talhelper/pkg/parse"
)

var (
	gencommandUpgradeOutDir     string
	gencommandUpgradeCfgFile    string
	gencommandUpgradeEnvFile    []string

	gencommandUpgradeFlagNode             string
	gencommandUpgradeExtraFlags           []string
	gencommandUpgradeInstallerRegistryURL string
)

var gencommandUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Generate talosctl upgrade commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := parse.ParseConfig(gencommandUpgradeCfgFile, gencommandUpgradeEnvFile)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		err = generate.GenerateUpgradeCommand(&cfg, gencommandUpgradeOutDir, gencommandUpgradeFlagNode, gencommandUpgradeInstallerRegistryURL, gencommandUpgradeExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl upgrade command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandUpgradeCmd)
	gencommandUpgradeCmd.Flags().StringVarP(&gencommandUpgradeCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	gencommandUpgradeCmd.Flags().StringVarP(&gencommandUpgradeOutDir, "out-dir", "o", "./clusterconfig", "Directory where the generated files were dumped with `genconfig`.")
	gencommandUpgradeCmd.Flags().StringSliceVar(&gencommandUpgradeEnvFile, "env-file", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	gencommandUpgradeCmd.Flags().StringSliceVar(&gencommandUpgradeExtraFlags, "extra-flags", []string{}, "List of additional flags that will be injected into the generated commands.")
	gencommandUpgradeCmd.Flags().StringVarP(&gencommandUpgradeFlagNode, "node", "n", "", "A specific node to generate the command for. If not specified, will generate for all nodes.")
	gencommandUpgradeCmd.Flags().StringVarP(&gencommandUpgradeInstallerRegistryURL, "registry-url", "r", "factory.talos.dev/installer", "Registry url of the image")
}
