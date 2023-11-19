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
	gencommandUpgradeOutDir     string
	gencommandUpgradeCfgFile    string
	gencommandUpgradeEnvFile    []string

	gencommandUpgradeFlagNode             string
	gencommandUpgradeExtraFlags           []string
	gencommandUpgradeInstallerRegistryURL string
)

var gencommandUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Generate talosctl upgrade commands.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfgByte, err := os.ReadFile(gencommandUpgradeCfgFile)
		if err != nil {
			log.Fatalf("failed to read config file: %s", err)
		}

		if err := substitute.LoadEnvFromFiles(gencommandUpgradeEnvFile); err != nil {
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

		err = generate.GenerateUpgradeCommand(&cfg, gencommandUpgradeOutDir, gencommandUpgradeFlagNode, gencommandUpgradeInstallerRegistryURL, gencommandUpgradeExtraFlags)
		if err != nil {
			log.Fatalf("failed to generate talosctl command: %s", err)
		}
	},
}

func init() {
	gencommandCmd.AddCommand(gencommandUpgradeCmd)
	gencommandUpgradeCmd.Flags().StringVarP(&gencommandUpgradeCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	gencommandUpgradeCmd.Flags().StringVarP(&gencommandUpgradeOutDir, "out-dir", "o", "./clusterconfig", "Directory where the generated files were dumped with `genconfig`.")
	gencommandUpgradeCmd.Flags().StringSliceVar(&gencommandUpgradeEnvFile, "env-file", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	gencommandUpgradeCmd.Flags().StringSliceVar(&gencommandUpgradeExtraFlags, "extra-flags", []string{""}, "List of additional flags that will be injected into the generated commands.")
	gencommandUpgradeCmd.Flags().StringVarP(&gencommandUpgradeFlagNode, "node", "n", "", "A specific node to generate the command for. If not specified, will generate for all nodes.")
	gencommandUpgradeCmd.Flags().StringVarP(&gencommandUpgradeInstallerRegistryURL, "registry-url", "r", "factory.talos.dev/installer", "Registry url of the image")
}
