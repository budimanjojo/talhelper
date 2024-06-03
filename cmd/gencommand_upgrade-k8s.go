package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/generate"
)

var gencommandUpgradeK8sCmd = &cobra.Command{
	Use:   "upgrade-k8s",
	Short: "Generate talosctl upgrade-k8s commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadAndValidateFromFile(gencommandCfgFile, gencommandEnvFile, false)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		err = generate.GenerateUpgradeK8sCommand(cfg, gencommandOutDir, gencommandNode, gencommandExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl upgrade-k8s command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandUpgradeK8sCmd)
}
