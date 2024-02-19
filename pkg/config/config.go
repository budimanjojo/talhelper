package config

import (
	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/siderolabs/talos/pkg/machinery/config/types/network"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/siderolabs/talos/pkg/machinery/nethelpers"
)

type TalhelperConfig struct {
	ClusterName                    string              `yaml:"clusterName" jsonschema:"required,description=Name of the cluster"`
	TalosVersion                   string              `yaml:"talosVersion,omitempty" jsonschema:"example=v1.5.4,description=Talos version to perform installation"`
	KubernetesVersion              string              `yaml:"kubernetesVersion,omitempty" jsonschema:"example=v1.27.0,description=Kubernetes version to use"`
	Endpoint                       string              `yaml:"endpoint" jsonschema:"required,example=https://192.168.200.10:6443,description=Cluster's controlplane endpoint"`
	Domain                         string              `yaml:"domain,omitempty" jsonschema:"example=cluster.local,description=The domain to be used by Kubernetes DNS"`
	AllowSchedulingOnMasters       bool                `yaml:"allowSchedulingOnMasters,omitempty" jsonschema:"description=Whether to allow running workload on controlplane nodes"`
	AllowSchedulingOnControlPlanes bool                `yaml:"allowSchedulingOnControlPlanes,omitempty" jsonschema:"description=Whether to allow running workload on controlplane nodes. It is an alias to \"AllowSchedulingOnMasters\""`
	AdditionalMachineCertSans      []string            `yaml:"additionalMachineCertSans,omitempty" jsonschema:"description=Extra certificate SANs for the machine's certificate"`
	AdditionalApiServerCertSans    []string            `yaml:"additionalApiServerCertSans,omitempty" jsonschema:"description=Extra certificate SANs for the API server's certificate"`
	ClusterPodNets                 []string            `yaml:"clusterPodNets,omitempty" jsonschema:"description=The pod subnet CIDR list"`
	ClusterSvcNets                 []string            `yaml:"clusterSvcNets,omitempty" jsonschema:"description=The service subnet CIDR list"`
	CNIConfig                      *v1alpha1.CNIConfig `yaml:"cniConfig,omitempty" jsonschema:"description=The CNI to be used for the cluster's network"`
	Patches                        []string            `yaml:"patches,omitempty" jsonschema:"description=Patches to be applied to all nodes"`
	Nodes                          []Node              `yaml:"nodes" jsonschema:"required,description=List of configurations for Node"`
	ImageFactory                   ImageFactory        `yaml:"imageFactory,omitempty" jsonschema:"Configuration for image factory"`
	ControlPlane                   NodeConfigs         `yaml:"controlPlane,omitempty" jsonschema:"description=Configurations targetted for all controlplane nodes"`
	Worker                         NodeConfigs         `yaml:"worker,omitempty" jsonschema:"description=Configurations targetted for all worker nodes"`
}

type Node struct {
	Hostname               string                        `yaml:"hostname" jsonschema:"required,description=Hostname of the node"`
	IPAddress              string                        `yaml:"ipAddress,omitempty" jsonschema:"required,example=192.168.200.11,description=IP address where the node can be reached, can also be a comma separated IP addresses"`
	ControlPlane           bool                          `yaml:"controlPlane" jsonschema:"description=Whether the node is a controlplane"`
	InstallDisk            string                        `yaml:"installDisk,omitempty" jsonschema:"oneof_required=installDiskSelector,description=The disk used for installation"`
	InstallDiskSelector    *v1alpha1.InstallDiskSelector `yaml:"installDiskSelector,omitempty" jsonschema:"oneof_required=installDisk,description=Look up disk used for installation"`
	IgnoreHostname         bool                          `yaml:"ignoreHostname" jsonschema:"description=Whether to set \"machine.network.hostname\" to the generated config file"`
	OverridePatches        bool                          `yaml:"overridePatches,omitempty" jsonschema:"description=Whether \"patches\" defined here should override the one defined in node group"`
	OverrideExtraManifests bool                          `yaml:"overrideExtraManifests,omitempty" jsonschema:"description=Whether \"extraManifests\" defined here should override the one defined in node group"`
	NodeConfigs            `yaml:",inline" jsonschema:"description=Node specific configurations that will override node group configurations"`
}

type NodeConfigs struct {
	NodeLabels          map[string]string              `yaml:"nodeLabels" jsonschema:"description=Labels to be added to the node"`
	NodeTaints          map[string]string              `yaml:"nodeTaints" jsonschema:"description=Node taints for the node. Effect is optional"`
	MachineDisks        []*v1alpha1.MachineDisk        `yaml:"machineDisks,omitempty" jsonschema:"description=List of additional disks to partition, format, mount"`
	MachineFiles        []*v1alpha1.MachineFile        `yaml:"machineFiles,omitempty" jsonschema:"description=List of files to create inside the node"`
	DisableSearchDomain bool                           `yaml:"disableSearchDomain,omitempty" jsonschema:"description=Whether to disable generating default search domain"`
	KernelModules       []*v1alpha1.KernelModuleConfig `yaml:"kernelModules,omitempty" jsonschema:"description=List of additional kernel modules to load inside the node"`
	Nameservers         []string                       `yaml:"nameservers,omitempty" jsonschema:"description=List of nameservers for the node"`
	NetworkInterfaces   []*v1alpha1.Device             `yaml:"networkInterfaces,omitempty" jsonschema:"description=List of network interface configuration for the node"`
	ExtraManifests      []string                       `yaml:"extraManifests,omitempty" jsonschema:"description=List of manifest files to be added to the node"`
	Patches             []string                       `yaml:"patches,omitempty" jsonschema:"description=Patches to be applied to the node"`
	TalosImageURL       string                         `yaml:"talosImageURL" jsonschema:"example=factory.talos.dev/installer/e9c7ef96884d4fbc8c0a1304ccca4bb0287d766a8b4125997cb9dbe84262144e,description=Talos installer image url for the node"`
	Schematic           *schematic.Schematic           `yaml:"schematic,omitempty" jsonschema:"description=Talos image customization to be used in the installer image"`
	MachineSpec         MachineSpec                    `yaml:"machineSpec,omitempty" jsonschema:"description=Machine hardware specification"`
	IngressFirewall     *IngressFirewall               `yaml:"ingressFirewall,omitempty" jsonschema:"description=Machine firewall specification"`
}

type ImageFactory struct {
	RegistryURL       string `yaml:"registryURL,omitempty" jsonschema:"default=factory.talos.dev,description=Registry url or the image"`
	SchematicEndpoint string `yaml:"schematicEndpoint,omitempty" jsonschema:"default=/schematics,description:Endpoint to get schematic ID from the registry"`
	Protocol          string `yaml:"protocol,omitempty" jsonschema:"default=https,description=Protocol of the registry(https or http)"`
	InstallerURLTmpl  string `yaml:"installerURLTmpl,omitempty" jsonschema:"default={{.RegistryURL}}/installer{{if .Secureboot}}-secureboot{{end}}/{{.ID}}:{{.Version}},description=Template for installer image URL"`
	ISOURLTmpl        string `yaml:"ISOURLTmpl,omitempty" jsonschema:"default={{.Protocol}}://{{.RegistryURL}}/image/{{.ID}}/{{.Version}}/{{.Mode}}-{{.Arch}}{{if .Secureboot}}-secureboot{{end}}{{if and .Secureboot .UseUKI}}-uki.efi{{else}}.iso{{end}},description=Template for ISO image URL"`
}

type MachineSpec struct {
	Mode       string `yaml:"mode,omitempty" jsonschema:"default=metal,description=Machine mode (e.g: metal)"`
	Arch       string `yaml:"arch,omitempty" jsonschema:"default=amd64,description=Machine architecture (e.g: amd64, arm64)"`
	Secureboot bool   `yaml:"secureboot,omitempty" jsonschema:"default=false,description=Whether to enable Secure Boot"`
	UseUKI     bool   `yaml:"useUKI,omitempty" jsonschema:"default=false,description=Whether to use UKI if Secure Boot is enabled"`
}

type IngressFirewall struct {
	DefaultAction nethelpers.DefaultAction `yaml:"defaultAction,omitempty" jsonschema:"default=block,description=Default action for all not explicitly configured traffic"`
	NetworkRules  []NetworkRule            `yaml:"rules,omitempty" jsonschema:"description=List of matching network rules to allow or block against the defaultAction"`
}

type NetworkRule struct {
	Name         string                   `yaml:"name" jsonschema:"description=Name of the rule"`
	PortSelector network.RulePortSelector `yaml:"portSelector" jsonschema:"description=Ports and protocols on the host affected by the rule"`
	Ingress      network.IngressConfig    `yaml:"ingress" jsonschema:"description=List of source subnets allowed to access the host ports/protocols"`
}
