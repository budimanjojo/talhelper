package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/generate"
)

var (
	gencommandUpgradeK8sOutDir  string
	gencommandUpgradeK8sCfgFile string
	gencommandUpgradeK8sEnvFile []string

	gencommandUpgradeK8sFlagNode   string
	gencommandUpgradeK8sExtraFlags []string
)

var gencommandUpgradeK8sCmd = &cobra.Command{
	Use:   "upgrade-k8s",
	Short: "Generate talosctl upgrade-k8s commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadAndValidateFromFile(gencommandUpgradeK8sCfgFile, gencommandUpgradeK8sEnvFile)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		err = generate.GenerateUpgradeK8sCommand(cfg, gencommandUpgradeK8sOutDir, gencommandUpgradeK8sFlagNode, gencommandUpgradeK8sExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl upgrade-k8s command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandUpgradeK8sCmd)
	gencommandUpgradeK8sCmd.Flags().StringVarP(&gencommandUpgradeK8sCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	gencommandUpgradeK8sCmd.Flags().StringVarP(&gencommandUpgradeK8sOutDir, "out-dir", "o", "./clusterconfig", "Directory where the generated files were dumped with `genconfig`.")
	gencommandUpgradeK8sCmd.Flags().StringSliceVar(&gencommandUpgradeK8sEnvFile, "env-file", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	gencommandUpgradeK8sCmd.Flags().StringSliceVar(&gencommandUpgradeK8sExtraFlags, "extra-flags", []string{}, "List of additional flags that will be injected into the generated commands.")
	gencommandUpgradeK8sCmd.Flags().StringVarP(&gencommandUpgradeK8sFlagNode, "node", "n", "", "A specific node to generate the command for. If not specified, will generate for all nodes.")
}
