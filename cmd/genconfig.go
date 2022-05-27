package cmd

import (
	"fmt"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/spf13/cobra"

	"sigs.k8s.io/yaml"
)

var (
	outDir     string
	configFile string
	varsFile   string
)

var (
	genconfigCmd = &cobra.Command{
		Use:   "genconfig",
		Short: "Generate Talos cluster config YAML file",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
			fmt.Println("genconfig called")
			fmt.Printf("configFile is %v, outDir is %v, varsFile is %v\n", configFile, outDir, varsFile)

			data, err := config.DecryptYamlWithSops(configFile)
			if err != nil {
				fmt.Println(err)
			}

			var m config.TalhelperConfig

			yaml.Unmarshal(data, &m)

			fmt.Println("Controlplane patches are: ", m.ControlPlane.ConfigPatches)

			err = m.GenerateConfig(outDir)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(genconfigCmd)

	genconfigCmd.Flags().StringVarP(&outDir, "out-dir", "o", "./clusterconfig", "Directory where to dump the generated files")
	genconfigCmd.Flags().StringVarP(&configFile, "config-file", "c", "config.yaml", "File containing configurations for nodes")
	genconfigCmd.Flags().StringVarP(&varsFile, "vars-file", "f", "vars.yaml", "File containing variables to load")
}
