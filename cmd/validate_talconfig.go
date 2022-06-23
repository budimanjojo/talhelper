package cmd

import (
	"fmt"
	"log"

	"github.com/budimanjojo/talhelper/pkg/config/validate"
	"github.com/spf13/cobra"
)

var (
	validateTHCmd = &cobra.Command{
		Use:   "talconfig [file]",
		Short: "Check the validity of talhelper config file",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := "talconfig.yaml"

			if len(args) > 0 {
				cfg = args[0]
			}

			found, err := validate.ValidateFromFile(cfg)
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

}
