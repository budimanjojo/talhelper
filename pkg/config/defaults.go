package config

import "github.com/talos-systems/talos/pkg/machinery/constants"

var (
	// renovate: depName=siderolabs/talos datasource=github-releases
	latestTalosVersion = "v1.0.6"
)

func (c TalhelperConfig) k8sVersion() string {
	if c.KubernetesVersion == "" {
		return constants.DefaultKubernetesVersion
	}
	return c.KubernetesVersion
}

func (c TalhelperConfig) talosVersion() string {
	if c.TalosVersion == "" {
		return latestTalosVersion
	}
	return c.TalosVersion
}
