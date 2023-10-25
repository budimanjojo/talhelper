package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	genurlISOCfgFile     string
	genurlISORegistryURL string
	genurlISOVersion     string
	genurlISOTalosMode   string
	genurlISOArch        string
	genurlISOExtensions  []string
	genurlISOKernelArgs  []string
)

var genurlISOCmd = &cobra.Command{
	Use:   "iso",
	Short: "Generate URL for Talos ISO image",
	Run: func(cmd *cobra.Command, args []string) {
		var version string
		if _, err := os.Stat(genurlISOCfgFile); err == nil {
			var m config.TalhelperConfig

			data, err := os.ReadFile(genurlISOCfgFile)
			if err != nil {
				log.Fatalf("failed to read Talhelper config file %s, %v", genurlISOCfgFile, err)
			}
			err = yaml.Unmarshal(data, &m)
			if err != nil {
				log.Fatalf("failed to unmarshal data: %s", err)
			}
			if m.TalosVersion != "" {
				version = m.TalosVersion
			} else {
				version = genurlISOVersion
			}
		} else if errors.Is(err, os.ErrNotExist) {
			version = genurlISOVersion
		} else {
			log.Fatalf("Failed to read Talhelper config file %s, %v", genurlISOCfgFile, err)
		}

		cfg := &schematic.Schematic{
			Customization: schematic.Customization{
				ExtraKernelArgs: genurlISOKernelArgs,
				SystemExtensions: schematic.SystemExtensions{
					OfficialExtensions: genurlISOExtensions,
				},
			},
		}

		url, err := talos.GetISOURL(cfg, genurlISORegistryURL, version, genurlISOTalosMode, genurlISOArch)
		if err != nil {
			log.Fatalf("Failed to generate installer url, %v", err)
		}

		fmt.Println(url)
	},
}

func init() {
	genurlCmd.AddCommand(genurlISOCmd)

	genurlISOCmd.Flags().StringVarP(&genurlISOCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	genurlISOCmd.Flags().StringVarP(&genurlISORegistryURL, "registry-url", "r", "https://factory.talos.dev/image", "Registry url of the image")
	genurlISOCmd.Flags().StringVarP(&genurlISOVersion, "version", "v", config.LatestTalosVersion, "Talos version to generate (defaults to latest Talos version)")
	genurlISOCmd.Flags().StringVarP(&genurlISOTalosMode, "talos-mode", "m", "metal", "Talos runtime mode to generate URL")
	genurlISOCmd.Flags().StringVarP(&genurlISOArch, "arch", "a", "amd64", "CPU architecture support of the image")
	genurlISOCmd.Flags().StringSliceVarP(&genurlISOExtensions, "extension", "e", []string{}, "Official extension image to be included in the image")
	genurlISOCmd.Flags().StringSliceVarP(&genurlISOKernelArgs, "kernel-arg", "k", []string{}, "Kernel arguments to be passed to the image kernel")
}
