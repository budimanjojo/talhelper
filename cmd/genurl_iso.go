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
	genurlISOTalosMode string
	genurlISOArch      string
	genurlISOUseUKI    bool
)

var genurlISOCmd = &cobra.Command{
	Use:   "iso",
	Short: "Generate URL for Talos ISO image",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(genurlCfgFile); err == nil {
			cfg, err := config.LoadAndValidateFromFile(genurlCfgFile, genurlEnvFile, false)
			if err != nil {
				log.Fatalf("failed to parse config file: %s", err)
			}

			var urls []string
			for _, node := range cfg.Nodes {
				if genurlNode != "" && !node.ContainsIP(genurlNode) && node.Hostname != genurlNode {
					continue
				}

				schema := &schematic.Schematic{}
				if node.Schematic != nil {
					schema = node.Schematic
				}

				if node.ContainsIP(genurlNode) || node.Hostname == genurlNode {
					url, err := talos.GetISOURL(schema, cfg.GetImageFactory(), node.GetMachineSpec(), cfg.GetTalosVersion(), genurlOfflineMode)
					if err != nil {
						log.Fatalf("Failed to generate ISO url for %s, %v", node.Hostname, err)
					}
					urls = append(urls, fmt.Sprintf(node.Hostname+": "+url))
					break
				}

				url, err := talos.GetISOURL(schema, cfg.GetImageFactory(), node.GetMachineSpec(), cfg.GetTalosVersion(), genurlOfflineMode)
				if err != nil {
					log.Fatalf("Failed to generate ISO url for %s, %v", node.Hostname, err)
				}
				urls = append(urls, fmt.Sprintf(node.Hostname+": "+url))
			}

			switch len(urls) {
			case 0:
				log.Fatalf("Node with IP address or hostname of %s is not found in the config file", genurlNode)
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
					ExtraKernelArgs: genurlKernelArgs,
					SystemExtensions: schematic.SystemExtensions{
						OfficialExtensions: genurlExtensions,
					},
				},
			}
			tcfg := &config.TalhelperConfig{}
			spec := &config.MachineSpec{
				Mode:       genurlISOTalosMode,
				Arch:       genurlISOArch,
				Secureboot: genurlSecureboot,
				UseUKI:     genurlISOUseUKI,
			}
			url, err := talos.GetISOURL(cfg, tcfg.GetImageFactory(), spec, genurlVersion, genurlOfflineMode)
			if err != nil {
				log.Fatalf("Failed to generate installer url, %v", err)
			}

			fmt.Println(url)
		} else {
			log.Fatalf("Failed to read Talhelper config file %s, %v", genurlCfgFile, err)
		}
	},
}

func init() {
	genurlCmd.AddCommand(genurlISOCmd)

	genurlISOCmd.Flags().StringVarP(&genurlISOTalosMode, "talos-mode", "m", "metal", "Talos runtime mode to generate URL")
	genurlISOCmd.Flags().StringVarP(&genurlISOArch, "arch", "a", "amd64", "CPU architecture support of the image")
	genurlISOCmd.Flags().BoolVar(&genurlISOUseUKI, "use-uki", false, "Whether to generate UKI image url if Secure Boot is enabled")
}
