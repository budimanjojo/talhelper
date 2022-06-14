package cmd

import (
	"log"

	"github.com/budimanjojo/talhelper/pkg/generate"
	"github.com/spf13/cobra"
)

var (
	gensecretPatchCfg bool
	gensecretFromCfg  string
	gensecretCfgFile  string
)

var gensecretCmd = &cobra.Command{
	Use:   "gensecret",
	Short: "Generate Talos cluster secrets",
	Run: func(cmd *cobra.Command, args []string) {
		err := generate.GenerateOutput(gensecretFromCfg)
		if err != nil {
			log.Fatalf("failed to generate secret bundle: %s", err)
		}

		if gensecretPatchCfg {
			err := generate.PatchTalhelperConfig(gensecretCfgFile)
			if err != nil {
				log.Fatalf("failed to patch talhelper config %s: %s", genconfigCfgFile, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(gensecretCmd)

	gensecretCmd.Flags().StringVarP(&gensecretCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	gensecretCmd.Flags().StringVarP(&gensecretFromCfg, "from-configfile", "f", "", "Talos cluster node configuration file to generate secret from")
	gensecretCmd.Flags().BoolVarP(&gensecretPatchCfg, "patch-configfile", "p", false, "Whether to generate inline patches into config file")
}
