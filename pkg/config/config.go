package config

type TalhelperConfig struct {
	ClusterName  string                 `yaml:"clusterName"`
	TalosVersion string                 `yaml:"talosVersion"`
	Endpoint     string                 `yaml:"endpoint"`
	Nodes        []nodes                `yaml:"nodes"`
	ControlPlane controlPlane           `yaml:"controlPlane"`
	Worker       worker                 `yaml:"worker"`
}

type nodes struct {
	Hostname     string `yaml:"hostname"`
	Domain       string `yaml:"domain"`
	IPAddress    string `yaml:"ipAddress"`
	ControlPlane bool   `yaml:"controlPlane"`
	InstallDisk  string `yaml:"installDisk"`
}

type controlPlane struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
}

type worker struct {
	ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
}
