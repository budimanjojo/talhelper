package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/generate"
)

var gencommandBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Generate talosctl bootstrap commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadAndValidateFromFile(gencommandCfgFile, gencommandEnvFile, false)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		err = generate.GenerateBootstrapCommand(cfg, gencommandOutDir, gencommandNode, gencommandExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl bootstrap command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandBootstrapCmd)
}
