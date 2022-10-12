package config

import (
	"strings"
	"net"

	tnet "github.com/talos-systems/net"
	"github.com/talos-systems/talos/pkg/machinery/constants"
)

var (
	// renovate: depName=siderolabs/talos datasource=github-releases
	latestTalosVersion = "v1.2.4"
)

func (c *TalhelperConfig) GetK8sVersion() string {
	if c.KubernetesVersion == "" {
		return ""
	}
	return strings.TrimPrefix(c.KubernetesVersion, "v")
}

func (c *TalhelperConfig) GetTalosVersion() string {
	if c.TalosVersion == "" {
		return latestTalosVersion
	}

	if !strings.HasPrefix(c.TalosVersion, "v") {
		return "v" + c.TalosVersion
	}
	return c.TalosVersion
}

func (c *TalhelperConfig) GetClusterPodNets() []string {
	if len(c.ClusterPodNets) == 0 {
		if tnet.IsIPv6(net.ParseIP(c.Endpoint)) {
			c.ClusterPodNets = []string{constants.DefaultIPv6PodNet}
		} else {
			c.ClusterPodNets = []string{constants.DefaultIPv4PodNet}
		}
	}
	return c.ClusterPodNets
}

func (c *TalhelperConfig) GetClusterSvcNets() []string {
	if len(c.ClusterSvcNets) == 0 {
		if tnet.IsIPv6(net.ParseIP(c.Endpoint)) {
			c.ClusterSvcNets = []string{constants.DefaultIPv6ServiceNet}
		} else {
			c.ClusterSvcNets = []string{constants.DefaultIPv4ServiceNet}
		}
	}
	return c.ClusterSvcNets
}

func (c *TalhelperConfig) GetInstallerURL() string {
	return "ghcr.io/siderolabs/installer:" + c.GetTalosVersion()
}
