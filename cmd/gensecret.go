package cmd

import (
	"log"

	"github.com/budimanjojo/talhelper/v3/pkg/generate"
	"github.com/spf13/cobra"
)

var gensecretFromCfg string

var GensecretCmd = &cobra.Command{
	Use:   "gensecret",
	Short: "Generate Talos cluster secrets",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := generate.GenerateSecret(gensecretFromCfg)
		if err != nil {
			log.Fatalf("failed to generate secret bundle: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(GensecretCmd)

	GensecretCmd.Flags().StringVarP(&gensecretFromCfg, "from-configfile", "f", "", "Talos cluster node configuration file to generate secret from")
}
