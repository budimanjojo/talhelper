package generate

import (
	"fmt"
	"log"
	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/talos"
)

func GenerateCommand(cfg *config.TalhelperConfig, gencommandFlagNode string, gencommandFlagApply bool, gencommandFlagUpgrade bool) error {
	if !gencommandFlagApply && !gencommandFlagUpgrade {
		return fmt.Errorf("Must select one of `--apply` of `--upgrade`\n")
	}

	for _, node := range cfg.Nodes {
		isSelectedNode := ( (gencommandFlagNode != "") && (gencommandFlagNode == node.IPAddress) )
		allNodesSelected := (gencommandFlagNode == "")

		if allNodesSelected || isSelectedNode {
			if gencommandFlagApply {
				fileName := outDir + "/" + cfg.ClusterName + "-" + node.Hostname + ".yaml"
				fmt.Printf("talosctl apply-config --talosconfig %s/talosconfig --nodes %s --file %s;\n", outDir, node.IPAddress, fileName)
			}
			
			if gencommandFlagUpgrade {
				if node.Schematic != nil {
					url, err := talos.GetInstallerURL(node.Schematic, gencommandInstallerRegistryURL, cfg.GetTalosVersion())
					if err != nil {
						return fmt.Errorf("Failed to generate installer url for %s, %v", node.Hostname, err)
					}
				} else {
					url, _ := talos.GetInstallerURL(&schematicfg.Schematic{}, gencommandInstallerRegistryURL, cfg.GetTalosVersion())
				}
				fmt.Printf("talosctl upgrade --talosconfig %s/talosconfig --nodes %s --image %s;\n", outDir, node.IPAddress, url)
			}
		}
	}

	return nil
}