package config

type TalhelperConfig struct {
	ClusterName       string       `yaml:"clusterName"`
	TalosVersion      string       `yaml:"talosVersion"`
	KubernetesVersion string       `yaml:"kubernetesVersion"`
	Endpoint          string       `yaml:"endpoint"`
	Nodes             []nodes      `yaml:"nodes"`
	CNIConfig         cniConfig    `yaml:"cniConfig"`
	ControlPlane      controlPlane `yaml:"controlPlane,omitempty"`
	Worker            worker       `yaml:"worker,omitempty"`
}

type nodes struct {
	Hostname     string                 `yaml:"hostname"`
	Domain       string                 `yaml:"domain"`
	IPAddress    string                 `yaml:"ipAddress"`
	ControlPlane bool                   `yaml:"controlPlane"`
	InstallDisk  string                 `yaml:"installDisk"`
	InlinePatch  map[string]interface{} `yaml:"inlinePatch"`
}

type cniConfig struct {
	Name string   `yaml:"name"`
	Urls []string `yaml:"urls"`
}

type controlPlane struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
	InlinePatch   map[string]interface{}   `yaml:"inlinePatch,omitempty"`
}

type worker struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
	InlinePatch   map[string]interface{}   `yaml:"inlinePatch,omitempty"`
}
