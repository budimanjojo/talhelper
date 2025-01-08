package cmd

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/talos"
	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/spf13/cobra"
)

var (
	genurlImageTalosMode  string
	genurlImageArch       string
	genurlImageUseUKI     bool
	genurlImageBootMethod string
	genurlImageSuffix     string
)

var genurlImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate URL for Talos ISO or disk image",
	Run: func(cmd *cobra.Command, args []string) {
		if !slices.Contains([]string{"iso", "disk-image", "pxe"}, genurlImageBootMethod) {
			log.Fatalf("invalid boot-method, should be one of iso, disk-image, pxe")
		}
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

				if node.ImageSchematic != nil {
					schema = node.ImageSchematic
				}

				if node.ContainsIP(genurlNode) || node.Hostname == genurlNode {
					url, err := talos.GetImageURL(schema, cfg.GetImageFactory(), node.GetMachineSpec(), cfg.GetTalosVersion(), genurlOfflineMode)
					if err != nil {
						log.Fatalf("Failed to generate ISO url for %s, %v", node.Hostname, err)
					}
					urls = append(urls, fmt.Sprintf("%s: %s", node.Hostname, url))
					break
				}

				url, err := talos.GetImageURL(schema, cfg.GetImageFactory(), node.GetMachineSpec(), cfg.GetTalosVersion(), genurlOfflineMode)
				if err != nil {
					log.Fatalf("Failed to generate ISO url for %s, %v", node.Hostname, err)
				}
				urls = append(urls, fmt.Sprintf("%s: %s", node.Hostname, url))
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
			slog.Debug("no config file found")
			slog.Debug("generating from provided flags", slog.Any("kernelArg", genurlKernelArgs), slog.Any("extension", genurlExtensions))
			cfg := &schematic.Schematic{
				Customization: schematic.Customization{
					ExtraKernelArgs: genurlKernelArgs,
					SystemExtensions: schematic.SystemExtensions{
						OfficialExtensions: genurlExtensions,
					},
				},
			}
			tcfg := &config.TalhelperConfig{}
			n := &config.Node{}
			n.MachineSpec.Mode = genurlImageTalosMode
			n.MachineSpec.Arch = genurlImageArch
			n.MachineSpec.Secureboot = genurlSecureboot
			n.MachineSpec.UseUKI = genurlImageUseUKI
			n.MachineSpec.BootMethod = genurlImageBootMethod
			n.MachineSpec.ImageSuffix = genurlImageSuffix

			url, err := talos.GetImageURL(cfg, tcfg.GetImageFactory(), &n.MachineSpec, genurlVersion, genurlOfflineMode)
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
	genurlCmd.AddCommand(genurlImageCmd)

	genurlImageCmd.Flags().StringVarP(&genurlImageTalosMode, "talos-mode", "m", "metal", "Talos runtime mode to generate URL")
	genurlImageCmd.Flags().StringVarP(&genurlImageArch, "arch", "a", "amd64", "CPU architecture support of the image")
	genurlImageCmd.Flags().BoolVar(&genurlImageUseUKI, "use-uki", false, "Whether to generate UKI image url if Secure Boot is enabled")
	genurlImageCmd.Flags().StringVar(&genurlImageBootMethod, "boot-method", "iso", "Boot method of the image (can be disk-image, iso, or pxe)")
	genurlImageCmd.Flags().StringVar(&genurlImageSuffix, "suffix", "", "The image file extension (only used when boot-method is not iso) (e.g: raw.xz, raw.tar.gz, qcow2)")
}
