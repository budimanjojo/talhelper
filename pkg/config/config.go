package config

type TalhelperConfig struct {
	ClusterName string `yaml:"clusterName"`
	TalosVersion string `yaml:"talosVersion"`
	Endpoint string `yaml:"endpoint"`
	Nodes []struct {
		Hostname string `yaml:"hostname"`
		Domain string `yaml:"domain"`
		IPAddress string `yaml:"ipAddress"`
		ControlPlane bool `yaml:"controlPlane"`
		InstallDisk string `yaml:"installDisk"`
	} `yaml:"nodes"`
	ControlPlane struct {
		ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
	} `yaml:"controlplane"`
	Worker struct {
		ConfigPatches []map[string]interface{} `yaml:"configPatches,omitempty"`
	} `yaml:"worker"`
}
