package cmd

import (
	"errors"
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
	genconfigOutDir      string
	genconfigCfgFile     string
	genconfigTalosMode   string
	genconfigNoGitignore bool
	genconfigEnvFile     []string
	genconfigSecretFile  []string
	genconfigDryRun      bool
)

var genconfigCmd = &cobra.Command{
	Use:   "genconfig",
	Short: "Generate Talos cluster config YAML files",
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

		var secretFile string
		for _, file := range genconfigSecretFile {
			if _, err := os.Stat(file); err == nil {
				secretFile = file
			} else if errors.Is(err, os.ErrNotExist) {
				continue
			} else {
				log.Fatalf("failed to stat secret file %s: %s ", file, err)
			}
		}

		err = generate.GenerateConfig(&cfg, genconfigDryRun, genconfigOutDir, secretFile, genconfigTalosMode)
		if err != nil {
			log.Fatalf("failed to generate talos config: %s", err)
		}

		if !genconfigNoGitignore && !genconfigDryRun {
			err = cfg.GenerateGitignore(genconfigOutDir)
			if err != nil {
				log.Fatalf("failed to generate gitignore file: %s", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(genconfigCmd)

	genconfigCmd.Flags().StringVarP(&genconfigOutDir, "out-dir", "o", "./clusterconfig", "Directory where to dump the generated files")
	genconfigCmd.Flags().StringVarP(&genconfigCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	genconfigCmd.Flags().StringSliceVarP(&genconfigEnvFile, "env-file", "e", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	genconfigCmd.Flags().StringSliceVarP(&genconfigSecretFile, "secret-file", "s", []string{"talsecret.yaml", "talsecret.sops.yaml", "talsecret.yml", "talsecret.sops.yml"}, "List of files containing secrets for the cluster")
	genconfigCmd.Flags().StringVarP(&genconfigTalosMode, "talos-mode", "m", "metal", "Talos runtime mode to validate generated config")
	genconfigCmd.Flags().BoolVar(&genconfigNoGitignore, "no-gitignore", false, "Create/update gitignore file too")
	genconfigCmd.Flags().BoolVarP(&genconfigDryRun, "dry-run", "n", false, "Skip generating manifests and show diff instead")
}
