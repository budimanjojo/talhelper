package config

import (
	"net/netip"
	"strings"

	"github.com/siderolabs/talos/pkg/machinery/constants"
)

// renovate: depName=siderolabs/talos datasource=github-releases
var LatestTalosVersion = "v1.5.5"

var OfficialExtensions = []string{
	"siderolabs/amd-ucode",
	"siderolabs/bnx2-bnx2x",
	"siderolabs/drbd",
	"siderolabs/gasket-driver",
	"siderolabs/gvisor",
	"siderolabs/hello-world-service",
	"siderolabs/i915-ucode",
	"siderolabs/intel-ucode",
	"siderolabs/iscsi-tools",
	"siderolabs/nut-client",
	"siderolabs/nvidia-container-toolkit",
	"siderolabs/nvidia-fabricmanager",
	"siderolabs/nvidia-open-gpu-kernel-modules",
	"siderolabs/qemu-guest-agent",
	"siderolabs/tailscale",
	"siderolabs/thunderbolt",
	"siderolabs/usb-modem-drivers",
	"siderolabs/zfs",
	"siderolabs/nonfree-kmod-nvidia",
}

// GetK8sVersion returns Kubernetes version string without `v` prefix.
func (c *TalhelperConfig) GetK8sVersion() string {
	if c.KubernetesVersion == "" {
		return ""
	}
	return strings.TrimPrefix(c.KubernetesVersion, "v")
}

// GetTalosVersion returns Talos version string prefixed with `v`.
func (c *TalhelperConfig) GetTalosVersion() string {
	if c.TalosVersion == "" {
		return LatestTalosVersion
	}

	if !strings.HasPrefix(c.TalosVersion, "v") {
		return "v" + c.TalosVersion
	}
	return c.TalosVersion
}

// GetClusterPodNets returns `ClusterPodNets` strings.
func (c *TalhelperConfig) GetClusterPodNets() []string {
	if len(c.ClusterPodNets) == 0 {
		if endpointisIPv6(c.Endpoint) {
			c.ClusterPodNets = []string{constants.DefaultIPv6PodNet}
		} else {
			c.ClusterPodNets = []string{constants.DefaultIPv4PodNet}
		}
	}
	return c.ClusterPodNets
}

// GetClusterSvcNets returns `ClusterSvcNets` strings.
func (c *TalhelperConfig) GetClusterSvcNets() []string {
	if len(c.ClusterSvcNets) == 0 {
		if endpointisIPv6(c.Endpoint) {
			c.ClusterSvcNets = []string{constants.DefaultIPv6ServiceNet}
		} else {
			c.ClusterSvcNets = []string{constants.DefaultIPv4ServiceNet}
		}
	}
	return c.ClusterSvcNets
}

// GetInstallerURL returns installer URL string.
func (c *TalhelperConfig) GetInstallerURL() string {
	if c.TalosImageURL != "" {
		return c.TalosImageURL + ":" + c.GetTalosVersion()
	}

	return "ghcr.io/siderolabs/installer:" + c.GetTalosVersion()
}

// GetImageFactory returns default `imageFactory` if not specified.
func (c *TalhelperConfig) GetImageFactory() *ImageFactory {
	result := ImageFactory{
		RegistryURL:       "factory.talos.dev",
		SchematicEndpoint: "/schematics",
		Protocol:          "https",
		InstallerURLTmpl:  "{{.RegistryURL}}/installer/{{.ID}}:{{.Version}}",
	}
	if c.ImageFactory.RegistryURL != "" {
		result.RegistryURL = c.ImageFactory.RegistryURL
	}
	if c.ImageFactory.SchematicEndpoint != "" {
		result.SchematicEndpoint = c.ImageFactory.SchematicEndpoint
	}
	if c.ImageFactory.Protocol != "" {
		result.Protocol = c.ImageFactory.Protocol
	}
	if c.ImageFactory.InstallerURLTmpl != "" {
		result.InstallerURLTmpl = c.ImageFactory.InstallerURLTmpl
	}
	return &result
}

// endpointisIPv6 returns true if string is IPv6 address.
func endpointisIPv6(ep string) bool {
	addr, err := netip.ParseAddr(ep)
	if err == nil && addr.Is6() {
		return true
	}
	return false
}
