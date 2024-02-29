package cmd

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/generate"
)

var (
	genconfigOutDir      string
	genconfigCfgFile     string
	genconfigTalosMode   string
	genconfigNoGitignore bool
	genconfigEnvFile     []string
	genconfigSecretFile  []string
	genconfigDryRun      bool
	genconfigOfflineMode bool
)

var genconfigCmd = &cobra.Command{
	Use:   "genconfig",
	Short: "Generate Talos cluster config YAML files",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadAndValidateFromFile(genconfigCfgFile, genconfigEnvFile)
		if err != nil {
			log.Fatalf("failed to parse config file: %s", err)
		}

		var secretFile string
		for _, file := range genconfigSecretFile {
			if _, err := os.Stat(file); err == nil {
				secretFile = file
				slog.Debug(fmt.Sprintf("secret file is set to %s", secretFile))
			} else if errors.Is(err, os.ErrNotExist) {
				continue
			} else {
				log.Fatalf("failed to stat secret file %s: %s ", file, err)
			}
		}

		slog.Debug("start generating config file")
		err = generate.GenerateConfig(cfg, genconfigDryRun, genconfigOutDir, secretFile, genconfigTalosMode, genconfigOfflineMode)
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
	genconfigCmd.Flags().BoolVar(&genconfigOfflineMode, "offline-mode", false, "Generate schematic ID without doing POST request to image-factory")
}
