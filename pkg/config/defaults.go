package config

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/netip"
	"slices"
	"strings"

	"github.com/budimanjojo/talhelper/v3/pkg/config/schemas/versiontags"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

// renovate: depName=siderolabs/talos datasource=github-releases
var LatestTalosVersion = "v1.10.5"

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
		InstallerURLTmpl:  "{{.RegistryURL}}/{{.Mode}}-installer{{if .Secureboot}}-secureboot{{end}}/{{.ID}}:{{.Version}}",
		ImageURLTmpl:      "{{.Protocol}}://{{.RegistryURL}}/image/{{.ID}}/{{.Version}}/{{.Mode}}-{{.Arch}}{{if .Secureboot}}-secureboot{{end}}{{if and .Secureboot .UseUKI}}-uki.efi{{ else }}{{.Suffix}}{{end}}",
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
	if c.ImageFactory.ImageURLTmpl != "" {
		result.ImageURLTmpl = c.ImageFactory.ImageURLTmpl
	}
	return result
}

// GetMachineSpec returns default `MachineSpec` for `Node` if not specified.
func (n *Node) GetMachineSpec() *MachineSpec {
	result := &MachineSpec{
		Mode:        "metal",
		Arch:        "amd64",
		BootMethod:  "iso",
		ImageSuffix: "iso",
	}
	if n.MachineSpec.BootMethod != "" {
		result.BootMethod = n.MachineSpec.BootMethod
	}
	if n.MachineSpec.Mode != "" {
		result.Mode = n.MachineSpec.Mode
	}
	if n.MachineSpec.Arch != "" {
		result.Arch = n.MachineSpec.Arch
	}
	if n.MachineSpec.ImageSuffix != "" {
		result.ImageSuffix = n.MachineSpec.ImageSuffix
	}
	result.Secureboot = n.MachineSpec.Secureboot
	result.UseUKI = n.MachineSpec.UseUKI
	return result
}

// GetIPAddresses returns list of IPaddresses
func (n *Node) GetIPAddresses() []string {
	var result []string
	ips := strings.Split(n.IPAddress, ",")
	for _, ip := range ips {
		result = append(result, strings.TrimSpace(ip))
	}
	return result
}

func (n *Node) GetFilenameTmpl() string {
	tmpl := "{{ .ClusterName }}-{{ .Hostname }}.yaml"
	if n.FilenameTmpl != "" {
		tmpl = n.FilenameTmpl
	}

	return tmpl
}

// ContainsIP returns true if `n.IPAddress` contains `ip`
func (n *Node) ContainsIP(ip string) bool {
	return slices.Contains(n.GetIPAddresses(), ip)
}

// endpointisIPv6 returns true if string is IPv6 address.
func endpointisIPv6(ep string) bool {
	addr, err := netip.ParseAddr(ep)
	if err == nil && addr.Is6() {
		return true
	}
	return false
}
