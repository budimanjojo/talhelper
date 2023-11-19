package generate

import (
	"fmt"
	"strings"
	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/siderolabs/image-factory/pkg/schematic"
)

func GenerateCommand(cfg *config.TalhelperConfig, gencommandOutDir string, gencommandFlagNode string, gencommandFlagApply bool, gencommandFlagUpgrade bool, gencommandInstallerRegistryURL string) error {
	if !gencommandFlagApply && !gencommandFlagUpgrade {
		return fmt.Errorf("Must select one of `--apply` of `--upgrade`\n")
	}

	for _, node := range cfg.Nodes {
		isSelectedNode := ( (gencommandFlagNode != "") && (gencommandFlagNode == node.IPAddress) )
		allNodesSelected := (gencommandFlagNode == "")

		if allNodesSelected || isSelectedNode {
			if gencommandFlagApply {
				applyFlags := []string{
					"--talosconfig=" + gencommandOutDir + "/talosconfig",
					"--nodes=" + node.IPAddress,
					"--file=" + gencommandOutDir + "/" + cfg.ClusterName + "-" + node.Hostname + ".yaml",
				}
				fmt.Printf("talosctl apply-config %s;\n", strings.Join(applyFlags, " "))
			}
			
			if gencommandFlagUpgrade {
				var imageUrl string

				if node.Schematic != nil {
					url, err := talos.GetInstallerURL(node.Schematic, gencommandInstallerRegistryURL, cfg.GetTalosVersion())
					if err != nil {
						return fmt.Errorf("Failed to generate installer url for %s, %v", node.Hostname, err)
					}
					imageUrl = url
				} else {
					url, _ := talos.GetInstallerURL(&schematic.Schematic{}, gencommandInstallerRegistryURL, cfg.GetTalosVersion())
					imageUrl = url
				}

				upgradeFlags := []string{
					"--talosconfig=" + gencommandOutDir + "/talosconfig",
					"--nodes=" + node.IPAddress,
					"--image=" + imageUrl,
				}
				fmt.Printf("talosctl upgrade %s;\n", strings.Join(upgradeFlags, " "))
			}
		}
	}

	return nil
}