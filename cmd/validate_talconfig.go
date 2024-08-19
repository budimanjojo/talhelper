package cmd

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/substitute"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	validateTHEnvFile      []string
	validateTHNoSubstitute bool
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

		slog.Debug("start loading and validating config file")
		slog.Debug(fmt.Sprintf("reading %s", cfg))
		cfgByte, err := config.FromFile(cfg)
		if err != nil {
			log.Fatalf("failed to read config file: %s", err)
		}

		if !validateTHNoSubstitute {
			if err := substitute.LoadEnvFromFiles(validateTHEnvFile); err != nil {
				log.Fatalf("failed to load env file: %s", err)
			}
			cfgByte, err = substitute.SubstituteEnvFromByte(cfgByte)
			if err != nil {
				log.Fatalf("failed trying to substitute env: %s", err)
			}
		}

		errs, warns, err := config.ValidateFromByte(cfgByte)
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
					fmt.Println(l)
				}
			}
		} else {
			fmt.Println("Your talhelper config file is looking great!")
		}
	},
}

func init() {
	validateCmd.AddCommand(validateTHCmd)

	validateTHCmd.Flags().StringSliceVarP(&validateTHEnvFile, "env-file", "e", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	validateTHCmd.Flags().BoolVar(&validateTHNoSubstitute, "no-substitute", false, "Whether to do envsubst on before validation")
}
