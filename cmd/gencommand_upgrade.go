package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/generate"
)

var gencommandUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Generate talosctl upgrade commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadAndValidateFromFile(gencommandCfgFile, gencommandEnvFile)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		err = generate.GenerateUpgradeCommand(cfg, gencommandOutDir, gencommandNode, gencommandExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl upgrade command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandUpgradeCmd)
}
