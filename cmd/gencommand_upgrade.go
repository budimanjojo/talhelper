package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/generate"
)

var gencommandOfflineMode bool

var gencommandUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Generate talosctl upgrade commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadAndValidateFromFile(gencommandCfgFile, gencommandEnvFile, false)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		err = generate.GenerateUpgradeCommand(cfg, gencommandOutDir, gencommandNode, gencommandExtraFlags, gencommandOfflineMode)
		if err != nil {
			log.Fatalf("failed to generate talosctl upgrade command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandUpgradeCmd)

	gencommandUpgradeCmd.Flags().BoolVar(&gencommandOfflineMode, "offline-mode", false, "Generate schematic ID without doing POST request to image-factory")
}
