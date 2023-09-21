package config

import "github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"

type TalhelperConfig struct {
	ClusterName                    string       `yaml:"clusterName"`
	TalosImageURL                  string       `yaml:"talosImageURL"`
	TalosVersion                   string       `yaml:"talosVersion,omitempty"`
	KubernetesVersion              string       `yaml:"kubernetesVersion,omitempty"`
	Endpoint                       string       `yaml:"endpoint"`
	Domain                         string       `yaml:"domain,omitempty"`
	AllowSchedulingOnMasters       bool         `yaml:"allowSchedulingOnMasters,omitempty"`
	AllowSchedulingOnControlPlanes bool         `yaml:"allowSchedulingOnControlPlanes,omitempty"`
	AdditionalMachineCertSans      []string     `yaml:"additionalMachineCertSans,omitempty"`
	AdditionalApiServerCertSans    []string     `yaml:"additionalApiServerCertSans,omitempty"`
	ClusterPodNets                 []string     `yaml:"clusterPodNets,omitempty"`
	ClusterSvcNets                 []string     `yaml:"clusterSvcNets,omitempty"`
	CNIConfig                      cniConfig    `yaml:"cniConfig,omitempty"`
	Patches                        []string     `yaml:"patches,omitempty"`
	Nodes                          []Node       `yaml:"nodes"`
	ControlPlane                   controlPlane `yaml:"controlPlane,omitempty"`
	Worker                         worker       `yaml:"worker,omitempty"`
}

type Node struct {
	Hostname            string                            `yaml:"hostname"`
	IPAddress           string                            `yaml:"ipAddress,omitempty"`
	ControlPlane        bool                              `yaml:"controlPlane"`
	NodeLabels          map[string]string                 `yaml:"nodeLabels"`
	InstallDisk         string                            `yaml:"installDisk,omitempty"`
	InstallDiskSelector *v1alpha1.InstallDiskSelector     `yaml:"installDiskSelector,omitempty"`
	MachineDisks        []*v1alpha1.MachineDisk           `yaml:"machineDisks,omitempty"`
	MachineFiles        []*v1alpha1.MachineFile           `yaml:"machineFiles,omitempty"`
	Extensions          []v1alpha1.InstallExtensionConfig `yaml:"extensions,omitempty"`
	DisableSearchDomain bool                              `yaml:"disableSearchDomain,omitempty"`
	KernelModules       []*v1alpha1.KernelModuleConfig    `yaml:"kernelModules,omitempty"`
	Nameservers         []string                          `yaml:"nameservers,omitempty"`
	NetworkInterfaces   []*v1alpha1.Device                `yaml:"networkInterfaces,omitempty"`
	ConfigPatches       []map[string]interface{}          `yaml:"configPatches,omitempty"`
	InlinePatch         map[string]interface{}            `yaml:"inlinePatch,omitempty"`
	Patches             []string                          `yaml:"patches,omitempty"`
}

type cniConfig struct {
	Name string   `yaml:"name"`
	Urls []string `yaml:"urls,omitempty"`
}

type controlPlane struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
	InlinePatch   map[string]interface{}   `yaml:"inlinePatch,omitempty"`
	Patches       []string                 `yaml:"patches,omitempty"`
}

type worker struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
	InlinePatch   map[string]interface{}   `yaml:"inlinePatch,omitempty"`
	Patches       []string                 `yaml:"patches,omitempty"`
}
