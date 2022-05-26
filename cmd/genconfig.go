package cmd

import (
	"fmt"
	"os"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/spf13/cobra"

	"sigs.k8s.io/yaml"
)

var (
	outDir     string
	configFile string
	varsFile   string
)

// type MachineConfig struct {
// 	Controlplane struct {
// 		Endpoint      string `yaml:"endpoint"`
// 		ConfigPatches patch.JsonPatch `yaml:"configPatches"`
// 	} `yaml:"controlplane"`
// }

var (
	genconfigCmd = &cobra.Command{
		Use:   "genconfig",
		Short: "Generate Talos cluster config YAML file",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
			fmt.Println("genconfig called")
			fmt.Printf("configFile is %v, outDir is %v, varsFile is %v\n", configFile, outDir, varsFile)

			data, err := os.ReadFile(configFile)
			if err != nil {
				fmt.Println(err)
			}

			var m config.TalhelperConfig

			yaml.Unmarshal(data, &m)
			// cpPatches, _ := json.Marshal(m.Controlplane.ConfigPatches)
			// fmt.Println(string(cpPatches))

			fmt.Println("Controlplane patches are: ", m.ControlPlane.ConfigPatches)

			// patch, _ := json.Marshal(m.ControlPlane.ConfigPatches)

			err = m.GenerateConfig(outDir)
			if err != nil {
				fmt.Println(err)
			}
			// marshaledCfgJson, _ := yaml.YAMLToJSON(marshaledCfg)
			//
			// decoded, _ := jsonpatch.DecodePatch(patch)
			// final, _ := decoded.Apply(marshaledCfgJson)
			// // fmt.Println(string(marshaledCfg))
			//
			// marshaledCfgYaml, _ := yaml.JSONToYAML(final)
			// fmt.Println(string(marshaledCfgYaml))
			//
			// clientCfg, _ := generate.Talosconfig(input, generate.WithEndpointList([]string{"172.0.0.1", "172.0.0.2"}))
			// marshaledClientCfg, _ := clientCfg.Bytes()
			// fmt.Println(string(marshaledClientCfg))
			// fmt.Println(out)
		},
	}
)

func init() {
	rootCmd.AddCommand(genconfigCmd)

	genconfigCmd.Flags().StringVarP(&outDir, "out-dir", "o", "./clusterconfig", "Directory where to dump the generated files")
	genconfigCmd.Flags().StringVarP(&configFile, "config-file", "c", "config.yaml", "File containing configurations for nodes")
	genconfigCmd.Flags().StringVarP(&varsFile, "vars-file", "f", "vars.yaml", "File containing variables to load")
}
