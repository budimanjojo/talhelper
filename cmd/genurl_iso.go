package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/substitute"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/fatih/color"
	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/spf13/cobra"
)

var (
	genurlISOCfgFile     string
	genurlISOEnvFile     []string
	genurlISONode        string
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
		if _, err := os.Stat(genurlISOCfgFile); err == nil {
			cfgByte, err := os.ReadFile(genurlISOCfgFile)
			if err != nil {
				log.Fatalf("failed to read config file: %s", err)
			}

			if err := substitute.LoadEnvFromFiles(genurlISOEnvFile); err != nil {
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

			var urls []string
			for _, node := range cfg.Nodes {
				if genurlISONode != "" && node.IPAddress != genurlISONode {
					continue
				}

				schema := &schematic.Schematic{}
				if node.Schematic != nil {
					schema = node.Schematic
				}

				if node.IPAddress == genurlISONode {
					url, err := talos.GetISOURL(schema, genurlISORegistryURL, cfg.GetTalosVersion(), genurlISOTalosMode, genurlISOArch)
					if err != nil {
						log.Fatalf("Failed to generate ISO url for %s, %v", node.Hostname, err)
					}
					urls = append(urls, fmt.Sprintf(node.Hostname+": "+url))
					break
				}

				url, err := talos.GetISOURL(schema, genurlISORegistryURL, cfg.GetTalosVersion(), genurlISOTalosMode, genurlISOArch)
				if err != nil {
					log.Fatalf("Failed to generate ISO url for %s, %v", node.Hostname, err)
				}
				urls = append(urls, fmt.Sprintf(node.Hostname+": "+url))
			}

			switch len(urls) {
			case 0:
				log.Fatalf("Node with IP address of %s is not found in the config file", genurlISONode)
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
					ExtraKernelArgs: genurlISOKernelArgs,
					SystemExtensions: schematic.SystemExtensions{
						OfficialExtensions: genurlISOExtensions,
					},
				},
			}
			url, err := talos.GetISOURL(cfg, genurlISORegistryURL, genurlISOVersion, genurlISOTalosMode, genurlISOArch)
			if err != nil {
				log.Fatalf("Failed to generate installer url, %v", err)
			}

			fmt.Println(url)
		} else {
			log.Fatalf("Failed to read Talhelper config file %s, %v", genurlISOCfgFile, err)
		}
	},
}

func init() {
	genurlCmd.AddCommand(genurlISOCmd)

	genurlISOCmd.Flags().StringVarP(&genurlISOCfgFile, "config-file", "c", "talconfig.yaml", "File containing configurations for talhelper")
	genurlISOCmd.Flags().StringSliceVar(&genurlISOEnvFile, "env-file", []string{"talenv.yaml", "talenv.sops.yaml", "talenv.yml", "talenv.sops.yml"}, "List of files containing env variables for config file")
	genurlISOCmd.Flags().StringVarP(&genurlISONode, "node", "n", "", "A specific node to generate command for. If not specified, will generate for all nodes (ignored when talconfig.yaml is not found)")
	genurlISOCmd.Flags().StringVarP(&genurlISORegistryURL, "registry-url", "r", "https://factory.talos.dev/image", "Registry url of the image")
	genurlISOCmd.Flags().StringVarP(&genurlISOVersion, "version", "v", config.LatestTalosVersion, "Talos version to generate (defaults to latest Talos version)")
	genurlISOCmd.Flags().StringVarP(&genurlISOTalosMode, "talos-mode", "m", "metal", "Talos runtime mode to generate URL")
	genurlISOCmd.Flags().StringVarP(&genurlISOArch, "arch", "a", "amd64", "CPU architecture support of the image")
	genurlISOCmd.Flags().StringSliceVarP(&genurlISOExtensions, "extension", "e", []string{}, "Official extension image to be included in the image (ignored when talconfig.yaml is found)")
	genurlISOCmd.Flags().StringSliceVarP(&genurlISOKernelArgs, "kernel-arg", "k", []string{}, "Kernel arguments to be passed to the image kernel (ignored when talconfig.yaml is found)")
}
