package cmd

import (
	"fmt"
	"log"

	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/spf13/cobra"
)

var validateNCTalosMode string

var validateNCCmd = &cobra.Command{
	Use:   "nodeconfig [file]",
	Short: "Check the validity of Talos node config file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("please specify at least 1 talos node config file you want to validate")
		}

		for _, arg := range args {
			err := talos.ValidateConfigFromFile(arg, validateNCTalosMode)
			if err != nil {
				log.Fatalf("failed to validate Talos node config file %s: %s", arg, err)
			} else {
				fmt.Printf("%s is valid for %s mode\n", arg, validateNCTalosMode)
			}
		}
	},
}

func init() {
	validateCmd.AddCommand(validateNCCmd)

	validateNCCmd.Flags().StringVarP(&validateNCTalosMode, "mode", "m", "metal", "Talos runtime mode to validate with")
}
