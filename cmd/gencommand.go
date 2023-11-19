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
	gencommandOutDir     string
	gencommandCfgFile    string
	gencommandEnvFile    []string

	gencommandFlagApply            bool
	gencommandFlagUpgrade          bool
	gencommandFlagNode             string
	gencommandInstallerRegistryURL string
)

var gencommandCmd = &cobra.Command{
	Use:   "gencommand",
	Short: "Generate talosctl commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfgByte, err := os.ReadFile(gencommandCfgFile)
		if err != nil {
			log.Fatalf("failed to read config file: %s", err)
		}

		if err := substitute.LoadEnvFromFiles(gencommandEnvFile); err != nil {
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

		err = generate.GenerateCommand(&cfg, gencommandOutDir, gencommandFlagNode, gencommandFlagApply, gencommandFlagUpgrade, gencommandInstallerRegistryURL)
		if err != nil {
			log.Fatalf("failed to generate talosctl command: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gencommandCmd)
	gencommandCmd.Flags().StringVarP(&gencommandCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	gencommandCmd.Flags().StringVarP(&gencommandOutDir, "out-dir", "o", "./clusterconfig", "Directory where the generated files were dumped with `genconfig`.")
	gencommandCmd.Flags().StringSliceVar(&gencommandEnvFile, "env-file", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	gencommandCmd.Flags().StringVarP(&gencommandFlagNode, "node", "n", "", "A specific node to generate the command for. If not specified, will generate for all nodes.")
	gencommandCmd.Flags().StringVarP(&gencommandInstallerRegistryURL, "registry-url", "r", "factory.talos.dev/installer", "Registry url of the image")
	gencommandCmd.Flags().BoolVarP(&gencommandFlagApply, "apply", "a", false, "Generate the talosctl apply commands.")
	gencommandCmd.Flags().BoolVarP(&gencommandFlagUpgrade, "upgrade", "u", false, "Generate the talosctl upgrade commands.")
}
