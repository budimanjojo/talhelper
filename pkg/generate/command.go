package generate

import (
	"fmt"
	"strings"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/talos"
	"github.com/siderolabs/image-factory/pkg/schematic"
)

// GenerateApplyCommand prints out `talosctl apply-config` command for selected node.
// `outDir` is directory where generated talosconfig and node manifest files are located.
// If `node` is empty string, it prints commands for all nodes in `cfg.Nodes`.
// It returns error, if any.
func GenerateApplyCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result []string
	for _, n := range cfg.Nodes {
		isSelectedByIP := ((node != "") && (n.ContainsIP(node)))
		isSelectedByHostname := ((node != "") && (node == n.Hostname))
		allNodesSelected := (node == "")

		if isSelectedByIP {
			filename, err := n.GetOutputFileName(cfg)
			if err != nil {
				return err
			}
			applyFlags := []string{
				"--talosconfig=" + outDir + "/talosconfig",
				"--nodes=" + node,
				"--file=" + outDir + "/" + filename,
			}
			applyFlags = append(applyFlags, extraFlags...)
			result = append(result, fmt.Sprintf("talosctl apply-config %s;", strings.Join(applyFlags, " ")))
		} else if allNodesSelected || isSelectedByHostname {
			for _, ip := range n.GetIPAddresses() {
				filename, err := n.GetOutputFileName(cfg)
				if err != nil {
					return err
				}
				applyFlags := []string{
					"--talosconfig=" + outDir + "/talosconfig",
					"--nodes=" + ip,
					"--file=" + outDir + "/" + filename,
				}
				applyFlags = append(applyFlags, extraFlags...)
				result = append(result, fmt.Sprintf("talosctl apply-config %s;", strings.Join(applyFlags, " ")))
			}
		}
	}

	if len(result) > 0 {
		for _, r := range result {
			fmt.Printf("%s\n", r)
		}
		return nil
	} else {
		return fmt.Errorf("node with IP or hostname %s not found", node)
	}
}

// GenerateUpgradeCommand prints out `talosctl upgrade` command for selected node.
// `outDir` is directory where talosconfig is located.
// If `node` is empty string, it prints commands for all nodes in `cfg.Nodes`.
// It returns error, if any.
func GenerateUpgradeCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string, offlineMode bool) error {
	var result []string
	for _, n := range cfg.Nodes {
		isSelectedByIP := ((node != "") && (n.ContainsIP(node)))
		isSelectedByHostname := ((node != "") && (node == n.Hostname))
		allNodesSelected := (node == "")

		var url string
		if n.TalosImageURL != "" {
			url = n.TalosImageURL + ":" + cfg.GetTalosVersion()
		} else if n.Schematic != nil {
			var err error
			url, err = talos.GetInstallerURL(n.Schematic, cfg.GetImageFactory(), n.GetMachineSpec(), cfg.GetTalosVersion(), offlineMode)
			if err != nil {
				return fmt.Errorf("failed to generate installer url for %s, %v", n.Hostname, err)
			}
		} else {
			url, _ = talos.GetInstallerURL(&schematic.Schematic{}, cfg.GetImageFactory(), n.GetMachineSpec(), cfg.GetTalosVersion(), offlineMode)
		}

		if isSelectedByIP {
			upgradeFlags := []string{
				"--talosconfig=" + outDir + "/talosconfig",
				"--nodes=" + node,
				"--image=" + url,
			}
			upgradeFlags = append(upgradeFlags, extraFlags...)
			result = append(result, fmt.Sprintf("talosctl upgrade %s;", strings.Join(upgradeFlags, " ")))
		} else if allNodesSelected || isSelectedByHostname {
			for _, ip := range n.GetIPAddresses() {
				upgradeFlags := []string{
					"--talosconfig=" + outDir + "/talosconfig",
					"--nodes=" + ip,
					"--image=" + url,
				}
				upgradeFlags = append(upgradeFlags, extraFlags...)
				result = append(result, fmt.Sprintf("talosctl upgrade %s;", strings.Join(upgradeFlags, " ")))
			}
		}
	}

	if len(result) > 0 {
		for _, r := range result {
			fmt.Printf("%s\n", r)
		}
		return nil
	} else {
		return fmt.Errorf("node with IP or hostname %s not found", node)
	}
}

// GenerateUpgradeK8sCommand prints out `talosctl upgrade-k8s` command for selected node.
// `outDir` is directory where talosconfig is located.
// If `node` is empty string, it prints command for the first controlplane node found
// in `cfg.Nodes`. It returns error if `node` is not found or is not controlplane.
func GenerateUpgradeK8sCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result string

	if cfg.KubernetesVersion == "" {
		return fmt.Errorf("`kubernetesVersion` is not defined in the configuration")
	}

	for _, n := range cfg.Nodes {
		isSelectedByIP := ((node != "") && (n.ContainsIP(node)))
		isSelectedByHostname := ((node != "") && (node == n.Hostname))
		noNodeSelected := (node == "")
		upgradeFlags := []string{
			"--talosconfig=" + outDir + "/talosconfig",
			"--to=v" + cfg.GetK8sVersion(),
		}

		if noNodeSelected && n.ControlPlane {
			upgradeFlags = append(upgradeFlags, extraFlags...)
			// Use the first IP address of the node
			upgradeFlags = append(upgradeFlags, "--nodes="+n.GetIPAddresses()[0])
			result = fmt.Sprintf("talosctl upgrade-k8s %s;", strings.Join(upgradeFlags, " "))
			break
		}

		if isSelectedByIP {
			if !n.ControlPlane {
				return fmt.Errorf("node with IP %s is not a controlplane node", node)
			}
			upgradeFlags = append(upgradeFlags, extraFlags...)
			upgradeFlags = append(upgradeFlags, "--nodes="+node)
			result = fmt.Sprintf("talosctl upgrade-k8s %s;", strings.Join(upgradeFlags, " "))
			break
		} else if isSelectedByHostname {
			if !n.ControlPlane {
				return fmt.Errorf("node with hostname %s is not a controlplane node", node)
			}
			upgradeFlags = append(upgradeFlags, extraFlags...)
			// Use the first IP address of the hostname
			upgradeFlags = append(upgradeFlags, "--nodes="+n.GetIPAddresses()[0])
			result = fmt.Sprintf("talosctl upgrade-k8s %s;", strings.Join(upgradeFlags, " "))
			break
		}
	}

	if result != "" {
		fmt.Printf("%s\n", result)
		return nil
	} else {
		return fmt.Errorf("node with IP or hostname %s not found", node)
	}
}

// GenerateBootstrapCommand prints out `talosctl bootstrap` command for selected node.
// `outDir` is directory where talosconfig is located.
// If `node` is empty string, it prints command for the first controlplane node found
// in `cfg.Nodes`. It returns error if `node` is not found or is not controlplane.
func GenerateBootstrapCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result string
	for _, n := range cfg.Nodes {
		isSelectedByIP := ((node != "") && (n.ContainsIP(node)))
		isSelectedByHostname := ((node != "") && (node == n.Hostname))
		noNodeSelected := (node == "")
		bootstrapFlags := []string{
			"--talosconfig=" + outDir + "/talosconfig",
		}
		if noNodeSelected && n.ControlPlane {
			bootstrapFlags = append(bootstrapFlags, extraFlags...)
			// Use the first IP address of the node
			bootstrapFlags = append(bootstrapFlags, "--nodes="+n.GetIPAddresses()[0])
			result = fmt.Sprintf("talosctl bootstrap %s;", strings.Join(bootstrapFlags, " "))
			break
		}
		if isSelectedByIP {
			if !n.ControlPlane {
				return fmt.Errorf("node with IP %s is not a controlplane node", node)
			}
			bootstrapFlags = append(bootstrapFlags, extraFlags...)
			bootstrapFlags = append(bootstrapFlags, "--nodes="+node)
			result = fmt.Sprintf("talosctl bootstrap %s;", strings.Join(bootstrapFlags, " "))
			break
		} else if isSelectedByHostname {
			if !n.ControlPlane {
				return fmt.Errorf("node with hostname %s is not a controlplane node", node)
			}
			bootstrapFlags = append(bootstrapFlags, extraFlags...)
			// Use the first IP address of the hostname
			bootstrapFlags = append(bootstrapFlags, "--nodes="+n.GetIPAddresses()[0])
			result = fmt.Sprintf("talosctl bootstrap %s;", strings.Join(bootstrapFlags, " "))
			break
		}
	}

	if result != "" {
		fmt.Printf("%s\n", result)
		return nil
	} else {
		return fmt.Errorf("node with IP or hostname %s not found", node)
	}
}

// GenerateKubeconfigCommand prints out `talosctl kubeconfig` command for selected node.
// `outDir` is directory where talosconfig is located.
// If `node` is empty string, it prints command for the first controlplane node found
// in `cfg.Nodes`. It returns error if `node` is not found or is not controlplane.
func GenerateKubeconfigCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result string
	for _, n := range cfg.Nodes {
		isSelectedByIP := ((node != "") && (n.ContainsIP(node)))
		isSelectedByHostname := ((node != "") && (node == n.Hostname))
		noNodeSelected := (node == "")
		kubeconfigFlags := []string{
			"--talosconfig=" + outDir + "/talosconfig",
		}
		if noNodeSelected && n.ControlPlane {
			kubeconfigFlags = append(kubeconfigFlags, extraFlags...)
			// Use the first IP address of the node
			kubeconfigFlags = append(kubeconfigFlags, "--nodes="+n.GetIPAddresses()[0])
			result = fmt.Sprintf("talosctl kubeconfig %s;", strings.Join(kubeconfigFlags, " "))
			break
		}
		if isSelectedByIP {
			if !n.ControlPlane {
				return fmt.Errorf("node with IP %s is not a controlplane node", node)
			}
			kubeconfigFlags = append(kubeconfigFlags, extraFlags...)
			kubeconfigFlags = append(kubeconfigFlags, "--nodes="+node)
			result = fmt.Sprintf("talosctl kubeconfig %s;", strings.Join(kubeconfigFlags, " "))
			break
		} else if isSelectedByHostname {
			if !n.ControlPlane {
				return fmt.Errorf("node with hostname %s is not a controlplane node", node)
			}
			kubeconfigFlags = append(kubeconfigFlags, extraFlags...)
			// Use the first IP address of the hostname
			kubeconfigFlags = append(kubeconfigFlags, "--nodes="+n.GetIPAddresses()[0])
			result = fmt.Sprintf("talosctl kubeconfig %s;", strings.Join(kubeconfigFlags, " "))
			break
		}
	}

	if result != "" {
		fmt.Printf("%s\n", result)
		return nil
	} else {
		return fmt.Errorf("node with IP or hostname %s not found", node)
	}
}

// GenarateResetCommand prints out `talosctl reset` command for selected node.
// `outDir` is directory where generated talosconfig and node manifest files are located.
// If `node` is empty string, it prints commands for all nodes in `cfg.Nodes`.
// It returns error, if any.
func GenerateResetCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result []string
	for _, n := range cfg.Nodes {
		isSelectedByIP := ((node != "") && (n.ContainsIP(node)))
		isSelectedByHostname := ((node != "") && (node == n.Hostname))
		allNodesSelected := (node == "")

		if isSelectedByIP {
			resetFlags := []string{
				"--talosconfig=" + outDir + "/talosconfig",
				"--nodes=" + node,
			}
			resetFlags = append(resetFlags, extraFlags...)
			result = append(result, fmt.Sprintf("talosctl reset %s;", strings.Join(resetFlags, " ")))
		} else if allNodesSelected || isSelectedByHostname {
			for _, ip := range n.GetIPAddresses() {
				resetFlags := []string{
					"--talosconfig=" + outDir + "/talosconfig",
					"--nodes=" + ip,
				}
				resetFlags = append(resetFlags, extraFlags...)
				result = append(result, fmt.Sprintf("talosctl reset %s;", strings.Join(resetFlags, " ")))
			}
		}
	}

	if len(result) > 0 {
		for _, r := range result {
			fmt.Printf("%s\n", r)
		}
		return nil
	} else {
		return fmt.Errorf("node with IP or hostname %s not found", node)
	}
}

func GenerateHealthCommand(cfg *config.TalhelperConfig, outDir string, node string, extraFlags []string) error {
	var result []string

	if node == "" {
		for _, n := range cfg.Nodes {
			if n.ControlPlane {
				healthFlags := []string{
					"--talosconfig=" + outDir + "/talosconfig",
					"--nodes=" + n.GetIPAddresses()[0],
				}
				healthFlags = append(healthFlags, extraFlags...)
				result = append(result, fmt.Sprintf("talosctl health %s;", strings.Join(healthFlags, " ")))
				break
			}
		}
	} else {
		for _, n := range cfg.Nodes {
			isSelectedByIP := (n.ContainsIP(node))
			isSelectedByHostname := (node == n.Hostname)

			if isSelectedByIP || isSelectedByHostname {
				healthFlags := []string{
					"--talosconfig=" + outDir + "/talosconfig",
					"--nodes=" + node,
				}
				healthFlags = append(healthFlags, extraFlags...)
				result = append(result, fmt.Sprintf("talosctl health %s;", strings.Join(healthFlags, " ")))
				break
			}
		}
	}
	if len(result) > 0 {
		for _, r := range result {
			fmt.Printf("%s\n", r)
		}
		return nil
	} else {
		return fmt.Errorf("node with IP or hostname %s not found", node)
	}
}
