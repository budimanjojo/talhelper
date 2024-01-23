package config

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/netip"
	"strings"
	"tsehelper/pkg/versiontags"

	"github.com/siderolabs/talos/pkg/machinery/constants"
)

// renovate: depName=siderolabs/talos datasource=github-releases
var LatestTalosVersion = "v1.6.2"

//go:embed schemas/talos-extensions.json
var schemaFile []byte

var OfficialExtensions = generateExtensionSchema(schemaFile)

func generateExtensionSchema(data []byte) versiontags.TalosVersionTags {
	var vTags versiontags.TalosVersionTags
	if err := json.Unmarshal(data, &vTags); err != nil {
		log.Fatal(err)
	}
	return vTags
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

// GetImageFactory returns default `imageFactory` if not specified.
func (c *TalhelperConfig) GetImageFactory() *ImageFactory {
	result := &ImageFactory{
		RegistryURL:       "factory.talos.dev",
		SchematicEndpoint: "/schematics",
		Protocol:          "https",
		InstallerURLTmpl:  "{{.RegistryURL}}/installer{{if .Secureboot}}-secureboot{{end}}/{{.ID}}:{{.Version}}",
		ISOURLTmpl:        "{{.Protocol}}://{{.RegistryURL}}/image/{{.ID}}/{{.Version}}/{{.Mode}}-{{.Arch}}{{if .Secureboot}}-secureboot{{end}}{{if and .Secureboot .UseUKI}}-uki.efi{{else}}.iso{{end}}",
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
	return result
}

// GetMachineSpec returns default `MachineSpec` for `Node` if not specified.
func (n *Node) GetMachineSpec() *MachineSpec {
	result := &MachineSpec{
		Mode: "metal",
		Arch: "amd64",
	}
	if n.MachineSpec.Mode != "" {
		result.Mode = n.MachineSpec.Mode
	}
	if n.MachineSpec.Arch != "" {
		result.Arch = n.MachineSpec.Arch
	}
	result.Secureboot = n.MachineSpec.Secureboot
	result.UseUKI = n.MachineSpec.UseUKI
	return result
}

// endpointisIPv6 returns true if string is IPv6 address.
func endpointisIPv6(ep string) bool {
	addr, err := netip.ParseAddr(ep)
	if err == nil && addr.Is6() {
		return true
	}
	return false
}
