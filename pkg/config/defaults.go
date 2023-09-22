package config

import (
	"net/netip"
	"strings"

	"github.com/siderolabs/talos/pkg/machinery/constants"
)

var (
	// renovate: depName=siderolabs/talos datasource=github-releases
	LatestTalosVersion = "v1.5.3"
)

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

// endpointisIPv6 returns true if string is IPv6 address.
func endpointisIPv6(ep string) bool {
	addr, err := netip.ParseAddr(ep)
	if err == nil && addr.Is6() {
		return true
	}
	return false
}
