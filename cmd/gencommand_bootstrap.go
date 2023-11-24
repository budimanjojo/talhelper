package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/generate"
)

var (
	gencommandBootstrapOutDir  string
	gencommandBootstrapCfgFile string
	gencommandBootstrapEnvFile []string

	gencommandBootstrapFlagNode   string
	gencommandBootstrapExtraFlags []string
)

var gencommandBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Generate talosctl bootstrap commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadAndValidateFromFile(gencommandBootstrapCfgFile, gencommandBootstrapEnvFile)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		err = generate.GenerateBootstrapCommand(cfg, gencommandBootstrapOutDir, gencommandBootstrapFlagNode, gencommandBootstrapExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl bootstrap command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandBootstrapCmd)
	gencommandBootstrapCmd.Flags().StringVarP(&gencommandBootstrapCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	gencommandBootstrapCmd.Flags().StringVarP(&gencommandBootstrapOutDir, "out-dir", "o", "./clusterconfig", "Directory where the generated files were dumped with `genconfig`.")
	gencommandBootstrapCmd.Flags().StringSliceVar(&gencommandBootstrapEnvFile, "env-file", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	gencommandBootstrapCmd.Flags().StringSliceVar(&gencommandBootstrapExtraFlags, "extra-flags", []string{}, "List of additional flags that will be injected into the generated commands.")
	gencommandBootstrapCmd.Flags().StringVarP(&gencommandBootstrapFlagNode, "node", "n", "", "A specific node to generate the command for. If not specified, will generate for all nodes.")
}
