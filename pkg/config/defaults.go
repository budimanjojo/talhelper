package config

import (
	"strings"

	"github.com/talos-systems/talos/pkg/machinery/constants"
)

var (
	// renovate: depName=siderolabs/talos datasource=github-releases
	latestTalosVersion = "v1.0.6"
)

func (c TalhelperConfig) k8sVersion() string {
	if c.KubernetesVersion == "" {
		return constants.DefaultKubernetesVersion
	}
	return strings.TrimPrefix(c.KubernetesVersion, "v")
}

func (c TalhelperConfig) talosVersion() string {
	if c.TalosVersion == "" {
		return latestTalosVersion
	}

	if !strings.HasPrefix(c.TalosVersion, "v") {
		return "v" + c.TalosVersion
	}
	return c.TalosVersion
}
