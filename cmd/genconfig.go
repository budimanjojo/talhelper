package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/config/validate"
	"github.com/budimanjojo/talhelper/pkg/decrypt"
	"github.com/budimanjojo/talhelper/pkg/generate"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
				if _, err := os.Stat(file); err == nil {
					env, err := decrypt.DecryptYamlWithSops(file)
					if err != nil {
						log.Fatalf("failed to decrypt/read env file %s: %s", file, err)
					}

					err = config.LoadEnv(env)
					if err != nil {
						log.Fatalf("failed to load env variables from file %s: %s", file, err)
					}
				} else if errors.Is(err, os.ErrNotExist) {
					continue
				} else {
					log.Fatalf("failed to stat env file %s: %s ", file, err)
				}
			}

			cfFile, err := config.SubstituteEnvFromYaml(cf)
			if err != nil {
				log.Fatalf("failed to substitute env: %s", err)
			}

			prob, err := validate.ValidateFromByte(cfFile)
			if err != nil {
				log.Fatalf("failed to validate talhelper config file: %s", err)
			}
			if prob != nil {
				fmt.Println("There are issues with your talhelper config file:")
				for _, v := range prob {
					fmt.Printf("- " + v.One() + "\n")
				}
				os.Exit(1)
			}

			var m config.TalhelperConfig

			err = yaml.Unmarshal(cfFile, &m)
			if err != nil {
				log.Fatalf("failed to unmarshal data: %s", err)
			}

			err = generate.GenerateConfig(&m, genconfigOutDir, genconfigTalosMode)
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
