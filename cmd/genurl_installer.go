package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/substitute"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/fatih/color"
	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/spf13/cobra"
)

var (
	genurlInstallerCfgFile     string
	genurlInstallerEnvFile     []string
	genurlInstallerRegistryURL string
	genurlInstallerVersion     string
	genurlInstallerExtensions  []string
	genurlInstallerKernelArgs  []string
)

var genurlInstallerCmd = &cobra.Command{
	Use:   "installer",
	Short: "Generate URL for Talos installer image",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(genurlInstallerCfgFile); err == nil {
			cfgByte, err := os.ReadFile(genurlInstallerCfgFile)
			if err != nil {
				log.Fatalf("failed to read config file: %s", err)
			}

			if err := substitute.LoadEnvFromFiles(genurlInstallerEnvFile); err != nil {
				log.Fatalf("failed to load env file: %s", err)
			}

			cfgByte, err = substitute.SubstituteEnvFromByte(cfgByte)
			if err != nil {
				log.Fatalf("failed to substitute env: %s", err)
			}

			cfg, err := config.NewFromByte(cfgByte)
			if err != nil {
				log.Fatalf("failed to unmarshal config file: %s", err)
			}

			errs, warns := cfg.Validate()
			if len(errs) > 0 || len(warns) > 0 {
				color.Red("There are issues with your talhelper config file:")
				grouped := make(map[string][]string)
				for _, v := range errs {
					grouped[v.Field] = append(grouped[v.Field], v.Message.Error())
				}
				for _, v := range warns {
					grouped[v.Field] = append(grouped[v.Field], v.Message)
				}
				for field, list := range grouped {
					color.Yellow("field: %q\n", field)
					for _, l := range list {
						fmt.Printf(l + "\n")
					}
				}
				if len(errs) > 0 {
					os.Exit(1)
				}
			}

			for _, node := range cfg.Nodes {
				if node.Schematic != nil {
					url, err := talos.GetInstallerURL(node.Schematic, genurlInstallerRegistryURL, cfg.GetTalosVersion())
					if err != nil {
						log.Fatalf("Failed to generate installer url for %s, %v", node.Hostname, err)
					}
					fmt.Printf("%s: %s\n", node.Hostname, url)
				} else {
					url, _ := talos.GetInstallerURL(&schematic.Schematic{}, genurlInstallerRegistryURL, cfg.GetTalosVersion())
					fmt.Printf("%s: %s\n", node.Hostname, url)
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

			url, err := talos.GetInstallerURL(cfg, genurlInstallerRegistryURL, genurlInstallerVersion)
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
	genurlInstallerCmd.Flags().StringVarP(&genurlInstallerRegistryURL, "registry-url", "r", "factory.talos.dev/installer", "Registry url of the image")
	genurlInstallerCmd.Flags().StringVarP(&genurlInstallerVersion, "version", "v", config.LatestTalosVersion, "Talos version to generate (defaults to latest Talos version)")
	genurlInstallerCmd.Flags().StringSliceVarP(&genurlInstallerExtensions, "extension", "e", []string{}, "Official extension image to be included in the image (ignored when talconfig.yaml is found)")
	genurlInstallerCmd.Flags().StringSliceVarP(&genurlInstallerKernelArgs, "kernel-arg", "k", []string{}, "Kernel arguments to be passed to the image kernel (ignored when talconfig.yaml is found)")
}
