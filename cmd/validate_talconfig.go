package cmd

import (
	"fmt"
	"log"

	"github.com/budimanjojo/talhelper/pkg/config/validate"
	"github.com/spf13/cobra"
)

var (
	validateTHCfgFile string
)

var (
	validateTHCmd = &cobra.Command{
		Use:   "talconfig [file]",
		Short: "Check the validity of talhelper config file",
		Run: func(cmd *cobra.Command, args []string) {
			found, err := validate.ValidateFromFile(validateTHCfgFile)
			if err != nil {
				log.Fatalf("failed to validate talhelper config file: %s", err)
			}
			if found != nil {
				fmt.Println("There are issues with your talhelper config file:")
				for _, v := range found {
					fmt.Printf("- " + v.One() + "\n")
				}
			} else {
				fmt.Println("Your talhelper config file is looking great!")
			}
		},
	}
)

func init() {
	validateCmd.AddCommand(validateTHCmd)

	validateTHCmd.Flags().StringVarP(&validateTHCfgFile, "config-file", "c", "talconfig.yaml", "Talhelper config file to validate")
}
