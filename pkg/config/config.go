package config

import "github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"

type TalhelperConfig struct {
	ClusterName              string       `yaml:"clusterName"`
	TalosVersion             string       `yaml:"talosVersion,omitempty"`
	KubernetesVersion        string       `yaml:"kubernetesVersion,omitempty"`
	Endpoint                 string       `yaml:"endpoint"`
	Domain                   string       `yaml:"domain,omitempty"`
	AllowSchedulingOnMasters bool         `yaml:"allowSchedulingOnMasters,omitempty"`
	ClusterPodNets           []string     `yaml:"clusterPodNets,omitempty"`
	ClusterSvcNets           []string     `yaml:"clusterSvcNets,omitempty"`
	CNIConfig                cniConfig    `yaml:"cniConfig,omitempty"`
	Nodes                    []Nodes      `yaml:"nodes"`
	ControlPlane             controlPlane `yaml:"controlPlane,omitempty"`
	Worker                   worker       `yaml:"worker,omitempty"`
}

type Nodes struct {
	Hostname          string                   `yaml:"hostname"`
	IPAddress         string                   `yaml:"ipAddress"`
	ControlPlane      bool                     `yaml:"controlPlane"`
	InstallDisk       string                   `yaml:"installDisk"`
	Nameservers       []string                 `yaml:"nameservers,omitempty"`
	NetworkInterfaces []*v1alpha1.Device       `yaml:"networkInterfaces,omitempty"`
	ConfigPatches     []map[string]interface{} `yaml:"configPatches,omitempty"`
	InlinePatch       map[string]interface{}   `yaml:"inlinePatch,omitempty"`
}

type cniConfig struct {
	Name string   `yaml:"name"`
	Urls []string `yaml:"urls,omitempty"`
}

type controlPlane struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
	InlinePatch   map[string]interface{}   `yaml:"inlinePatch,omitempty"`
}

type worker struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
	InlinePatch   map[string]interface{}   `yaml:"inlinePatch,omitempty"`
}

type route struct {
	Network string `yaml:"network,omitempty"`
	Gateway string `yaml:"gateway,omitempty"`
	Source  string `yaml:"source,omitempty"`
	Metric  uint32 `yaml:"metric,omitempty"`
}
