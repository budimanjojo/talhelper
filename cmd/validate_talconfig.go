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

		errs, warns, err := config.ValidateFromFile(cfg)
		if err != nil {
			log.Fatalf("failed to validate talhelper config file: %s", err)
		}

		if len(errs) > 0 || len(warns) > 0 {
			color.Red("There are issues with your talhelper config file:")
			grouped := make(map[string][]string)
			for _, v := range errs {
				grouped[v.Field] = append(grouped[v.Field], v.Message.Error())
			}
			for _, v := range warns {
				grouped[v.Field] = append(grouped[v.Field], v.Message)
			}
			for field, list := range grouped {
				color.Yellow("field: %q\n", field)
				for _, l := range list {
					fmt.Printf(l + "\n")
				}
			}
		} else {
			fmt.Println("Your talhelper config file is looking great!")
		}
	},
}

func init() {
	validateCmd.AddCommand(validateTHCmd)
}
