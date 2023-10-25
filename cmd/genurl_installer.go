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
	genurlInstallerCfgFile     string
	genurlInstallerRegistryURL string
	genurlInstallerVersion     string
	genurlInstallerExtensions  []string
	genurlInstallerKernelArgs  []string
)

var genurlInstallerCmd = &cobra.Command{
	Use:   "installer",
	Short: "Generate URL for Talos installer image",
	Run: func(cmd *cobra.Command, args []string) {
		var version string
		if _, err := os.Stat(genurlInstallerCfgFile); err == nil {
			var m config.TalhelperConfig

			err = yaml.Unmarshal([]byte(genurlInstallerCfgFile), &m)
			if err != nil {
				log.Fatalf("failed to unmarshal data: %s", err)
			}
			if m.TalosVersion != "" {
				version = m.TalosVersion
			} else {
				version = genurlInstallerVersion
			}
		} else if errors.Is(err, os.ErrNotExist) {
			version = genurlInstallerVersion
		} else {
			log.Fatalf("Failed to read Talhelper config file %s, %v", genurlInstallerCfgFile, err)
		}

		cfg := &schematic.Schematic{
			Customization: schematic.Customization{
				ExtraKernelArgs: genurlInstallerKernelArgs,
				SystemExtensions: schematic.SystemExtensions{
					OfficialExtensions: genurlInstallerExtensions,
				},
			},
		}

		url, err := talos.GetInstallerURL(cfg, genurlInstallerRegistryURL, version)
		if err != nil {
			log.Fatalf("Failed to generate installer url, %v", err)
		}

		fmt.Println(url)
	},
}

func init() {
	genurlCmd.AddCommand(genurlInstallerCmd)

	genurlInstallerCmd.Flags().StringVarP(&genurlInstallerCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	genurlInstallerCmd.Flags().StringVarP(&genurlInstallerRegistryURL, "registry-url", "r", "factory.talos.dev/installer", "Registry url of the image")
	genurlInstallerCmd.Flags().StringVarP(&genurlInstallerVersion, "version", "v", config.LatestTalosVersion, "Talos version to generate (defaults to latest Talos version)")
	genurlInstallerCmd.Flags().StringSliceVarP(&genurlInstallerExtensions, "extension", "e", []string{}, "Official extension image to be included in the image")
	genurlInstallerCmd.Flags().StringSliceVarP(&genurlInstallerKernelArgs, "kernel-arg", "k", []string{}, "Kernel arguments to be passed to the image kernel")
}
