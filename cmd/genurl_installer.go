package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/spf13/cobra"
)

var (
	genurlInstallerCfgFile     string
	genurlInstallerEnvFile     []string
	genurlInstallerNode        string
	genurlInstallerRegistryURL string
	genurlInstallerVersion     string
	genurlInstallerExtensions  []string
	genurlInstallerKernelArgs  []string
	genurlInstallerOfflineMode bool
)

var genurlInstallerCmd = &cobra.Command{
	Use:   "installer",
	Short: "Generate URL for Talos installer image",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(genurlInstallerCfgFile); err == nil {
			cfg, err := config.LoadAndValidateFromFile(genurlInstallerCfgFile, genurlInstallerEnvFile)
			if err != nil {
				log.Fatalf("failed to parse config file: %s", err)
			}

			var urls []string
			for _, node := range cfg.Nodes {
				if genurlInstallerNode != "" && node.IPAddress != genurlInstallerNode {
					continue
				}

				schema := &schematic.Schematic{}
				if node.Schematic != nil {
					schema = node.Schematic
				}

				if node.IPAddress == genurlInstallerNode {
					url, err := talos.GetInstallerURL(schema, cfg.GetImageFactory(), cfg.GetTalosVersion(), genurlInstallerOfflineMode)
					if err != nil {
						log.Fatalf("Failed to generate installer url for %s, %v", node.Hostname, err)
					}
					urls = append(urls, fmt.Sprintf(node.Hostname+": "+url))
					break
				}

				url, err := talos.GetInstallerURL(schema, cfg.GetImageFactory(), cfg.GetTalosVersion(), genurlInstallerOfflineMode)
				if err != nil {
					log.Fatalf("Failed to generate installer url for %s, %v", node.Hostname, err)
				}
				urls = append(urls, fmt.Sprintf(node.Hostname+": "+url))
			}

			switch len(urls) {
			case 0:
				log.Fatalf("Node with IP address of %s is not found in the config file", genurlInstallerNode)
			case 1:
				s := strings.Split(urls[0], " ")
				fmt.Printf("%s\n", s[len(s)-1])
			default:
				for _, v := range urls {
					fmt.Printf("%s\n", v)
				}
			}
		} else if errors.Is(err, os.ErrNotExist) {
			cfg := &schematic.Schematic{
				Customization: schematic.Customization{
					ExtraKernelArgs: genurlInstallerKernelArgs,
					SystemExtensions: schematic.SystemExtensions{
						OfficialExtensions: genurlInstallerExtensions,
					},
				},
			}
			tconfig := &config.TalhelperConfig{}
			tconfig.ImageFactory.RegistryURL = genurlInstallerRegistryURL

			url, err := talos.GetInstallerURL(cfg, tconfig.GetImageFactory(), genurlInstallerVersion, genurlInstallerOfflineMode)
			if err != nil {
				log.Fatalf("Failed to generate installer url, %v", err)
			}

			fmt.Println(url)
		} else {
			log.Fatalf("Failed to read Talhelper config file %s, %v", genurlInstallerCfgFile, err)
		}
	},
}

func init() {
	genurlCmd.AddCommand(genurlInstallerCmd)

	genurlInstallerCmd.Flags().StringVarP(&genurlInstallerCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	genurlInstallerCmd.Flags().StringSliceVar(&genurlInstallerEnvFile, "env-file", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	genurlInstallerCmd.Flags().StringVarP(&genurlInstallerNode, "node", "n", "", "A specific node to generate command for. If not specified, will generate for all nodes (ignored when talconfig.yaml is not found)")
	genurlInstallerCmd.Flags().StringVarP(&genurlInstallerRegistryURL, "registry-url", "r", "factory.talos.dev", "Registry url of the image")
	genurlInstallerCmd.Flags().StringVarP(&genurlInstallerVersion, "version", "v", config.LatestTalosVersion, "Talos version to generate (defaults to latest Talos version)")
	genurlInstallerCmd.Flags().StringSliceVarP(&genurlInstallerExtensions, "extension", "e", []string{}, "Official extension image to be included in the image (ignored when talconfig.yaml is found)")
	genurlInstallerCmd.Flags().StringSliceVarP(&genurlInstallerKernelArgs, "kernel-arg", "k", []string{}, "Kernel arguments to be passed to the image kernel (ignored when talconfig.yaml is found)")
	genurlInstallerCmd.Flags().BoolVar(&genurlInstallerOfflineMode, "offline-mode", false, "Generate schematic ID without doing POST request to image-factory")
}
