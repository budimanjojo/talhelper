package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/generate"
)

var gencommandApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Generate talosctl apply-config commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadAndValidateFromFile(gencommandCfgFile, gencommandEnvFile, false)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		err = generate.GenerateApplyCommand(cfg, gencommandOutDir, gencommandNode, gencommandExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl apply command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandApplyCmd)
}
