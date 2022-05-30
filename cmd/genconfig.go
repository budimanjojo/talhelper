package cmd

import (
	"log"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/spf13/cobra"

	"sigs.k8s.io/yaml"
)

var (
	outDir      string
	configFile  string
	noGitignore bool
)

var (
	genconfigCmd = &cobra.Command{
		Use:   "genconfig",
		Short: "Generate Talos cluster config YAML file",
		Run: func(cmd *cobra.Command, args []string) {
			data, err := config.DecryptYamlWithSops(configFile)
			if err != nil {
				log.Fatalf("failed to decrypt/read file: %s", err)
			}

			var m config.TalhelperConfig

			err = yaml.Unmarshal(data, &m)
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
	genconfigCmd.Flags().StringVarP(&configFile, "config-file", "c", "config.yaml", "File containing configurations for nodes")
	genconfigCmd.Flags().BoolVar(&noGitignore, "no-gitignore", false, "Create/update gitignore file too")
}
