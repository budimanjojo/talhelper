package config

type TalhelperConfig struct {
	ClusterName       string       `yaml:"clusterName"`
	TalosVersion      string       `yaml:"talosVersion,omitempty"`
	KubernetesVersion string       `yaml:"kubernetesVersion,omitempty"`
	Endpoint          string       `yaml:"endpoint"`
	Domain            string       `yaml:"domain"`
	Nodes             []nodes      `yaml:"nodes"`
	CNIConfig         cniConfig    `yaml:"cniConfig,omitempty"`
	ControlPlane      controlPlane `yaml:"controlPlane,omitempty"`
	Worker            worker       `yaml:"worker,omitempty"`
}

type nodes struct {
	Hostname      string                   `yaml:"hostname"`
	IPAddress     string                   `yaml:"ipAddress"`
	ControlPlane  bool                     `yaml:"controlPlane"`
	InstallDisk   string                   `yaml:"installDisk"`
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
	InlinePatch   map[string]interface{}   `yaml:"inlinePatch,omitempty"`
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
