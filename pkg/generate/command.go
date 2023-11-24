package generate

import (
	"fmt"
	"strings"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/siderolabs/image-factory/pkg/schematic"
)

func GenerateApplyCommand(cfg *config.TalhelperConfig, gencommandOutDir string, gencommandFlagNode string, gencommandExtraFlags []string) error {
	for _, node := range cfg.Nodes {
		isSelectedNode := ((gencommandFlagNode != "") && (gencommandFlagNode == node.IPAddress))
		allNodesSelected := (gencommandFlagNode == "")

		if allNodesSelected || isSelectedNode {
			applyFlags := []string{
				"--talosconfig=" + gencommandOutDir + "/talosconfig",
				"--nodes=" + node.IPAddress,
				"--file=" + gencommandOutDir + "/" + cfg.ClusterName + "-" + node.Hostname + ".yaml",
			}
			applyFlags = append(applyFlags, gencommandExtraFlags...)
			fmt.Printf("talosctl apply-config %s;\n", strings.Join(applyFlags, " "))
		}
	}

	return nil
}

func GenerateUpgradeCommand(cfg *config.TalhelperConfig, gencommandOutDir string, gencommandFlagNode string, gencommandInstallerRegistryURL string, gencommandExtraFlags []string) error {
	for _, node := range cfg.Nodes {
		isSelectedNode := ((gencommandFlagNode != "") && (gencommandFlagNode == node.IPAddress))
		allNodesSelected := (gencommandFlagNode == "")

		if allNodesSelected || isSelectedNode {
			var url string
			if node.Schematic != nil {
				var err error
				url, err = talos.GetInstallerURL(node.Schematic, gencommandInstallerRegistryURL, cfg.GetTalosVersion())
				if err != nil {
					return fmt.Errorf("Failed to generate installer url for %s, %v", node.Hostname, err)
				}
			} else {
				url, _ = talos.GetInstallerURL(&schematic.Schematic{}, gencommandInstallerRegistryURL, cfg.GetTalosVersion())
			}

			upgradeFlags := []string{
				"--talosconfig=" + gencommandOutDir + "/talosconfig",
				"--nodes=" + node.IPAddress,
				"--image=" + url,
			}
			upgradeFlags = append(upgradeFlags, gencommandExtraFlags...)
			fmt.Printf("talosctl upgrade %s;\n", strings.Join(upgradeFlags, " "))
		}
	}

	return nil
}

func GenerateBootstrapCommand(cfg *config.TalhelperConfig, gencommandOutDir string, gencommandFlagNode string, gencommandExtraFlags []string) error {
	for _, node := range cfg.Nodes {
		isSelectedNode := ((gencommandFlagNode != "") && (gencommandFlagNode == node.IPAddress))
		noNodeSelected := (gencommandFlagNode == "")
		bootstrapFlags := []string{
			"--talosconfig=" + gencommandOutDir + "/talosconfig",
		}
		if noNodeSelected && node.ControlPlane {
			bootstrapFlags = append(bootstrapFlags, gencommandExtraFlags...)
			bootstrapFlags = append(bootstrapFlags, "--nodes="+node.IPAddress)
			fmt.Printf("talosctl bootstrap %s;\n", strings.Join(bootstrapFlags, " "))
			break
		}
		if isSelectedNode {
			if !node.ControlPlane {
				return fmt.Errorf("%s is not a controlplane node", node.IPAddress)
			}
			bootstrapFlags = append(bootstrapFlags, gencommandExtraFlags...)
			bootstrapFlags = append(bootstrapFlags, "--nodes="+node.IPAddress)
			fmt.Printf("talosctl bootstrap %s;\n", strings.Join(bootstrapFlags, " "))
			break
		}
	}

	return nil
}
