package generate

import (
	"fmt"
	"strings"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/siderolabs/image-factory/pkg/schematic"
)

// GenerateApplyCommand prints out `talosctl apply-config` command for selected node.
// `outDir` is directory where generated talosconfig and node manifest files are located.
// If `node` is empty string, it prints commands for all nodes in `cfg.Nodes`.
// It returns error, if any.
func GenerateApplyCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result []string
	for _, n := range cfg.Nodes {
		isSelectedNode := ((node != "") && (node == n.IPAddress))
		allNodesSelected := (node == "")

		if allNodesSelected || isSelectedNode {
			applyFlags := []string{
				"--talosconfig=" + outDir + "/talosconfig",
				"--nodes=" + n.IPAddress,
				"--file=" + outDir + "/" + cfg.ClusterName + "-" + n.Hostname + ".yaml",
			}
			applyFlags = append(applyFlags, extraFlags...)
			result = append(result, fmt.Sprintf("talosctl apply-config %s;", strings.Join(applyFlags, " ")))
		}
	}

	if len(result) > 0 {
		for _, r := range result {
			fmt.Printf("%s\n", r)
		}
		return nil
	} else {
		return fmt.Errorf("node with IP %s not found", node)
	}
}

// GenerateUpgradeCommand prints out `talosctl upgrade` command for selected node.
// `outDir` is directory where talosconfig is located.
// If `node` is empty string, it prints commands for all nodes in `cfg.Nodes`.
// It returns error, if any.
func GenerateUpgradeCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result []string
	for _, n := range cfg.Nodes {
		isSelectedNode := ((node != "") && (node == n.IPAddress))
		allNodesSelected := (node == "")

		if allNodesSelected || isSelectedNode {
			var url string
			if n.Schematic != nil {
				var err error
				url, err = talos.GetInstallerURL(n.Schematic, cfg.GetImageFactory(), cfg.GetTalosVersion(), true)
				if err != nil {
					return fmt.Errorf("Failed to generate installer url for %s, %v", n.Hostname, err)
				}
			} else {
				url, _ = talos.GetInstallerURL(&schematic.Schematic{}, cfg.GetImageFactory(), cfg.GetTalosVersion(), true)
			}

			upgradeFlags := []string{
				"--talosconfig=" + outDir + "/talosconfig",
				"--nodes=" + n.IPAddress,
				"--image=" + url,
			}
			upgradeFlags = append(upgradeFlags, extraFlags...)
			result = append(result, fmt.Sprintf("talosctl upgrade %s;", strings.Join(upgradeFlags, " ")))
		}
	}

	if len(result) > 0 {
		for _, r := range result {
			fmt.Printf("%s\n", r)
		}
		return nil
	} else {
		return fmt.Errorf("node with IP %s not found", node)
	}
}

// GenerateUpgradeK8sCommand prints out `talosctl upgrade-k8s` command for selected node.
// `outDir` is directory where talosconfig is located.
// If `node` is empty string, it prints command for the first controlplane node found
// in `cfg.Nodes`. It returns error if `node` is not found or is not controlplane.
func GenerateUpgradeK8sCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result string
	for _, n := range cfg.Nodes {
		isSelectedNode := ((node != "") && (node == n.IPAddress))
		noNodeSelected := (node == "")
		upgradeFlags := []string{
			"--talosconfig=" + outDir + "/talosconfig",
			"--to=v" + cfg.GetK8sVersion(),
		}

		if noNodeSelected && n.ControlPlane {
			upgradeFlags = append(upgradeFlags, extraFlags...)
			upgradeFlags = append(upgradeFlags, "--nodes="+n.IPAddress)
			result = fmt.Sprintf("talosctl upgrade-k8s %s;", strings.Join(upgradeFlags, " "))
			break
		}
		if isSelectedNode {
			if !n.ControlPlane {
				return fmt.Errorf("node with IP %s is not a controlplane node", n.IPAddress)
			}
			upgradeFlags = append(upgradeFlags, extraFlags...)
			upgradeFlags = append(upgradeFlags, "--nodes="+n.IPAddress)
			result = fmt.Sprintf("talosctl upgrade-k8s %s;", strings.Join(upgradeFlags, " "))
			break
		}
	}

	if result != "" {
		fmt.Printf("%s\n", result)
		return nil
	} else {
		return fmt.Errorf("node with IP %s not found", node)
	}
}

// GenerateBootstrapCommand prints out `talosctl bootstrap` command for selected node.
// `outDir` is directory where talosconfig is located.
// If `node` is empty string, it prints command for the first controlplane node found
// in `cfg.Nodes`. It returns error if `node` is not found or is not controlplane.
func GenerateBootstrapCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result string
	for _, n := range cfg.Nodes {
		isSelectedNode := ((node != "") && (node == n.IPAddress))
		noNodeSelected := (node == "")
		bootstrapFlags := []string{
			"--talosconfig=" + outDir + "/talosconfig",
		}
		if noNodeSelected && n.ControlPlane {
			bootstrapFlags = append(bootstrapFlags, extraFlags...)
			bootstrapFlags = append(bootstrapFlags, "--nodes="+n.IPAddress)
			result = fmt.Sprintf("talosctl bootstrap %s;", strings.Join(bootstrapFlags, " "))
			break
		}
		if isSelectedNode {
			if !n.ControlPlane {
				return fmt.Errorf("node with IP %s is not a controlplane node", n.IPAddress)
			}
			bootstrapFlags = append(bootstrapFlags, extraFlags...)
			bootstrapFlags = append(bootstrapFlags, "--nodes="+n.IPAddress)
			result = fmt.Sprintf("talosctl bootstrap %s;", strings.Join(bootstrapFlags, " "))
			break
		}
	}

	if result != "" {
		fmt.Printf("%s\n", result)
		return nil
	} else {
		return fmt.Errorf("node with IP %s not found", node)
	}
}
