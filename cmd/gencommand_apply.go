package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/generate"
	"github.com/budimanjojo/talhelper/pkg/substitute"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

var (
	gencommandApplyOutDir     string
	gencommandApplyCfgFile    string
	gencommandApplyEnvFile    []string

	gencommandApplyNode             string
	gencommandApplyExtraFlags           []string
)

var gencommandApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Generate talosctl apply-config commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfgByte, err := os.ReadFile(gencommandApplyCfgFile)
		if err != nil {
			log.Fatalf("failed to read config file: %s", err)
		}

		if err := substitute.LoadEnvFromFiles(gencommandApplyEnvFile); err != nil {
			log.Fatalf("failed to load env file: %s", err)
		}

		cfgByte, err = substitute.SubstituteEnvFromByte(cfgByte)
		if err != nil {
			log.Fatalf("failed to substitute env: %s", err)
		}

		cfg, err := config.NewFromByte(cfgByte)
		if err != nil {
			log.Fatalf("failed to unmarshal config file: %s", err)
		}

		errs, warns := cfg.Validate()
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

			if len(errs) > 0 {
				os.Exit(1)
			}
		}

		err = generate.GenerateApplyCommand(&cfg, gencommandApplyOutDir, gencommandApplyNode, gencommandApplyExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl apply command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandApplyCmd)
	gencommandApplyCmd.Flags().StringVarP(&gencommandApplyCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	gencommandApplyCmd.Flags().StringVarP(&gencommandApplyOutDir, "out-dir", "o", "./clusterconfig", "Directory that contains the generated config files to apply.")
	gencommandApplyCmd.Flags().StringSliceVar(&gencommandApplyEnvFile, "env-file", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	gencommandApplyCmd.Flags().StringSliceVar(&gencommandApplyExtraFlags, "extra-flags", []string{}, "List of additional flags that will be injected into the generated commands.")
	gencommandApplyCmd.Flags().StringVarP(&gencommandApplyNode, "node", "n", "", "A specific node to generate the command for. If not specified, will generate for all nodes.")
}
