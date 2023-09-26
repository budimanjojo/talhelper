package cmd

import (
	"fmt"
	"log"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var validateTHCmd = &cobra.Command{
	Use:   "talconfig [file]",
	Short: "Check the validity of talhelper config file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := "talconfig.yaml"

		if len(args) > 0 {
			cfg = args[0]
		}

		found, err := config.ValidateFromFile(cfg)
		if err != nil {
			log.Fatalf("failed to validate talhelper config file: %s", err)
		}
		if len(found) > 0 {
			color.Red("There are issues with your talhelper config file:")
			for _, v := range found {
				color.Yellow("field: %q\n", v.Field)
				fmt.Printf(v.Message.Error() + "\n")
			}
		} else {
			fmt.Println("Your talhelper config file is looking great!")
		}
	},
}

func init() {
	validateCmd.AddCommand(validateTHCmd)
}
