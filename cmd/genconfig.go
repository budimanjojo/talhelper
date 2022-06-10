package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/spf13/cobra"

	"sigs.k8s.io/yaml"
)

var (
	genconfigOutDir      string
	genconfigCfgFile     string
	genconfigTalosMode   string
	genconfigNoGitignore bool
	genconfigEnvFile     []string
)

var (
	genconfigCmd = &cobra.Command{
		Use:   "genconfig",
		Short: "Generate Talos cluster config YAML file",
		Run: func(cmd *cobra.Command, args []string) {
			cf, err := os.ReadFile(genconfigCfgFile)
			if err != nil {
				log.Fatalf("failed to read config file: %s", err)
			}

			for _, file := range genconfigEnvFile {
				if _, err := os.Stat(file); errors.Is(err, os.ErrExist) {
					env, err := config.DecryptYamlWithSops(file)
					if err != nil {
						log.Fatalf("failed to decrypt/read env file %s: %s", file, err)
					}

					config.LoadEnv(env)
				}
			}

			cfFile, err := config.SubstituteEnvFromYaml(cf)
			if err != nil {
				log.Fatalf("failed to substitute env: %s", err)
			}

			var m config.TalhelperConfig

			err = yaml.Unmarshal(cfFile, &m)
			if err != nil {
				log.Fatalf("failed to unmarshal data: %s", err)
			}

			err = m.GenerateConfig(genconfigOutDir, genconfigTalosMode)
			if err != nil {
				log.Fatalf("failed to generate talos config: %s", err)
			}

			if !genconfigNoGitignore {
				err = m.GenerateGitignore(genconfigOutDir)
				if err != nil {
					log.Fatalf("failed to generate gitignore file: %s", err)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(genconfigCmd)

	genconfigCmd.Flags().StringVarP(&genconfigOutDir, "out-dir", "o", "./clusterconfig", "Directory where to dump the generated files")
	genconfigCmd.Flags().StringVarP(&genconfigCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	genconfigCmd.Flags().StringSliceVarP(&genconfigEnvFile, "env-file", "e", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	genconfigCmd.Flags().StringVarP(&genconfigTalosMode, "talos-mode", "m", "metal", "Talos runtime mode to validate generated config")
	genconfigCmd.Flags().BoolVar(&genconfigNoGitignore, "no-gitignore", false, "Create/update gitignore file too")
}
