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
	generateApply      bool
	generateUpgrade      bool
	generateForNode      string
)

var gencommandCmd = &cobra.Command{
	Use:   "gencommand",
	Short: "Generate the talosctl command for applying the talhelper generated config files.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfgByte, err := os.ReadFile(genconfigCfgFile)
		if err != nil {
			log.Fatalf("failed to read config file: %s", err)
		}

		if err := substitute.LoadEnvFromFiles(genconfigEnvFile); err != nil {
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

		err = generate.GenerateCommand(&cfg, genconfigOutDir, generateForNode, generateApply, generateUpgrade)
		if err != nil {
			log.Fatalf("failed to generate talosctl command: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gencommandCmd)
	gencommandCmd.Flags().StringVarP(&genconfigCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	gencommandCmd.Flags().StringVarP(&genconfigOutDir, "out-dir", "o", "./clusterconfig", "Directory where to dump the generated files")
	gencommandCmd.Flags().StringVarP(&generateForNode, "node", "n", "", "A specific node to generate the command for. If not specified, will generate for all nodes.")
	gencommandCmd.Flags().BoolVarP(&generateApply, "apply", "a", false, "Generate the talosctl apply commands.")
	gencommandCmd.Flags().BoolVarP(&generateUpgrade, "upgrade", "u", false, "Generate the talosctl upgrade commands.")
}
