package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/a8m/envsubst"
	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/spf13/cobra"

	"sigs.k8s.io/yaml"
)

var (
	genconfigCmd = &cobra.Command{
		Use:   "genconfig",
		Short: "Generate Talos cluster config YAML file",
		Run: func(cmd *cobra.Command, args []string) {
			cf, err := os.ReadFile(configFile)
			cfFile := &cf
			if err != nil {
				log.Fatalf("failed to decrypt/read file: %s", err)
			}

			if _, err := os.Stat(envFile); !errors.Is(err, os.ErrNotExist) {
				env, err := config.DecryptYamlWithSops(envFile)
				if err != nil {
					log.Fatalf("failed to decrypt/read env file: %s", err)
				}

				*cfFile, err = config.SubstituteEnvFromYaml(env, cf)
				if err != nil {
					log.Fatalf("failed to substite env: %s", err)
				}
			}

			_, err = envsubst.Bytes(*cfFile)
			if err != nil {
				log.Fatalf("failed to substite env: %s", err)
			}

			var m config.TalhelperConfig

			err = yaml.Unmarshal(*cfFile, &m)
			if err != nil {
				log.Fatalf("failed to unmarshal data: %s", err)
			}

			err = m.GenerateConfig(outDir)
			if err != nil {
				log.Fatalf("failed to generate talos config: %s", err)
			}

			if !noGitignore {
				err = m.GenerateGitignore(outDir)
				if err != nil {
					log.Fatalf("failed to generate gitignore file: %s", err)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(genconfigCmd)

	genconfigCmd.Flags().StringVarP(&outDir, "out-dir", "o", "./clusterconfig", "Directory where to dump the generated files")
	genconfigCmd.Flags().StringVarP(&configFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	genconfigCmd.Flags().StringVarP(&envFile, "env-file", "e", "talenv.yaml", "File containing env variables for config file")
	genconfigCmd.Flags().BoolVar(&noGitignore, "no-gitignore", false, "Create/update gitignore file too")
}
