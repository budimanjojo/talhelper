package config

import (
	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
)

type TalhelperConfig struct {
	ClusterName                    string       `yaml:"clusterName" jsonschema:"required,description=Name of the cluster"`
	TalosImageURL                  string       `yaml:"talosImageURL" jsonschema:"default=ghcr.io/siderolabs/installer,description=The image URL used to perform installation"`
	TalosVersion                   string       `yaml:"talosVersion,omitempty" jsonschema:"example=v1.5.4,description=Talos version to perform installation"`
	KubernetesVersion              string       `yaml:"kubernetesVersion,omitempty" jsonschema:"example=v1.27.0,description=Kubernetes version to use"`
	Endpoint                       string       `yaml:"endpoint" jsonschema:"required,example=https://192.168.200.10:6443,description=Cluster's controlplane endpoint"`
	Domain                         string       `yaml:"domain,omitempty" jsonschema:"example=cluster.local,description=The domain to be used by Kubernetes DNS"`
	AllowSchedulingOnMasters       bool         `yaml:"allowSchedulingOnMasters,omitempty" jsonschema:"description=Whether to allow running workload on controlplane nodes"`
	AllowSchedulingOnControlPlanes bool         `yaml:"allowSchedulingOnControlPlanes,omitempty" jsonschema:"description=Whether to allow running workload on controlplane nodes. It is an alias to \"AllowSchedulingOnMasters\""`
	AdditionalMachineCertSans      []string     `yaml:"additionalMachineCertSans,omitempty" jsonschema:"description=Extra certificate SANs for the machine's certificate"`
	AdditionalApiServerCertSans    []string     `yaml:"additionalApiServerCertSans,omitempty" jsonschema:"description=Extra certificate SANs for the API server's certificate"`
	ClusterPodNets                 []string     `yaml:"clusterPodNets,omitempty" jsonschema:"description=The pod subnet CIDR list"`
	ClusterSvcNets                 []string     `yaml:"clusterSvcNets,omitempty" jsonschema:"description=The service subnet CIDR list"`
	CNIConfig                      cniConfig    `yaml:"cniConfig,omitempty" jsonschema:"description=The CNI to be used for the cluster's network"`
	Patches                        []string     `yaml:"patches,omitempty" jsonschema:"description=ClusterInlineManifest"`
	Nodes                          []Node       `yaml:"nodes" jsonschema:"required,description=List of configurations for Node"`
	ControlPlane                   controlPlane `yaml:"controlPlane,omitempty" jsonschema:"description=Configurations targetted for controlplane nodes"`
	Worker                         worker       `yaml:"worker,omitempty" jsonschema:"description=Configurations targetted for worker nodes"`
}

type Node struct {
	Hostname            string                            `yaml:"hostname" jsonschema:"required,description=Hostname of the node"`
	IPAddress           string                            `yaml:"ipAddress,omitempty" jsonschema:"required,example=192.168.200.11,description=IP address where the node can be reached"`
	ControlPlane        bool                              `yaml:"controlPlane" jsonschema:"description=Whether the node is a controlplane"`
	NodeLabels          map[string]string                 `yaml:"nodeLabels" jsonschema:"description=Labels to be added to the node"`
	InstallDisk         string                            `yaml:"installDisk,omitempty" jsonschema:"oneof_required=installDiskSelector,description=The disk used for installation"`
	InstallDiskSelector *v1alpha1.InstallDiskSelector     `yaml:"installDiskSelector,omitempty" jsonschema:"oneof_required=installDisk,description=Look up disk used for installation"`
	MachineDisks        []*v1alpha1.MachineDisk           `yaml:"machineDisks,omitempty" jsonschema:"description=List of additional disks to partition, format, mount"`
	MachineFiles        []*v1alpha1.MachineFile           `yaml:"machineFiles,omitempty" jsonschema:"description=List of files to create inside the node"`
	Extensions          []v1alpha1.InstallExtensionConfig `yaml:"extensions,omitempty" jsonschema:"description=DEPRECATED, use \"schematic\" instead"`
	DisableSearchDomain bool                              `yaml:"disableSearchDomain,omitempty" jsonschema:"description=Whether to disable generating default search domain"`
	KernelModules       []*v1alpha1.KernelModuleConfig    `yaml:"kernelModules,omitempty" jsonschema:"description=List of additional kernel modules to load inside the node"`
	Nameservers         []string                          `yaml:"nameservers,omitempty" jsonschema:"description=List of nameservers for the node"`
	NetworkInterfaces   []*v1alpha1.Device                `yaml:"networkInterfaces,omitempty" jsonschema:"description=List of network interface configuration for the node"`
	ConfigPatches       []map[string]interface{}          `yaml:"configPatches,omitempty" jsonschema:"description=DEPRECATED, use \"patches\" instead"`
	InlinePatch         map[string]interface{}            `yaml:"inlinePatch,omitempty" jsonschema:"description=DEPRECATED, use \"patches\" instead"`
	Patches             []string                          `yaml:"patches,omitempty" jsonschema:"description=Patches to be applied to the node"`
	TalosImageURL       string                            `yaml:"talosImageURL" jsonschema:"example=factory.talos.dev/installer/e9c7ef96884d4fbc8c0a1304ccca4bb0287d766a8b4125997cb9dbe84262144e,description=Talos installer image url for the node"`
	Schematic           *schematic.Schematic              `yaml:"schematic,omitempty" jsonschema:"description=Talos image customization to be used in the installer image"`
}

type cniConfig struct {
	Name string   `yaml:"name" jsonschema:"required,description=The name of CNI to use"`
	Urls []string `yaml:"urls,omitempty" jsonschema:"description=List of URLs containing manifests to apply for the CNI"`
}

type controlPlane struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty" jsonschema:"description=DEPRECATED, use \"patches\" instead"`
	InlinePatch   map[string]interface{}   `yaml:"inlinePatch,omitempty" jsonschema:"description=DEPRECATED, use \"patches\" instead"`
	Patches       []string                 `yaml:"patches,omitempty" jsonschema:"description=Patches to be applied to all controlplane nodes"`
}

type worker struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty" jsonschema:"description=DEPRECATED, use \"patches\" instead"`
	InlinePatch   map[string]interface{}   `yaml:"inlinePatch,omitempty" jsonschema:"description=DEPRECATED, use \"patches\" instead"`
	Patches       []string                 `yaml:"patches,omitempty" jsonschema:"description=Patches to be applied to all worker nodes"`
}
