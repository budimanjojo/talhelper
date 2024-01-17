package config

type InstallDiskSelectorWrapper struct {
	Size     string `yaml:"size" jsonschema:"description=Disk size,example=4GB"`
	Name     string `yaml:"name" jsonschema:"Disk name"`
	Model    string `yaml:"model" jsonschema:"Disk model"`
	Serial   string `yaml:"serial" jsonschema:"Disk serial number"`
	Modalias string `yaml:"modalias" jsonschema:"Disk modalias"`
	UUID     string `yaml:"uuid" jsonschema:"Disk UUID"`
	WWID     string `yaml:"wwid" jsonschema:"Disk WWID"`
	Type     string `yaml:"type" jsonschema:"Disk type,example=ssd"`
	BusPath  string `yaml:"busPath" jsonschema:"Disk bus path"`
}

type IngressFirewallWrapper struct {
	DefaultAction string               `yaml:"defaultAction" jsonschema:"default=block,description=Default action for all not explicitly configured traffic"`
	NetworkRules  []NetworkRuleWrapper `yaml:"rules" jsonschema:"description=List of matching network rules to allow or block against the defaultAction"`
}

type NetworkRuleWrapper struct {
	Name         string                 `yaml:"name" jsonschema:"description=Name of the rule"`
	PortSelector PortSelectorWrapper    `yaml:"portSelector" jsonschema:"description=Ports and protocols on the host affected by the rule"`
	Ingress      []IngressConfigWrapper `yaml:"ingress" jsonschema:"description=List of source subnets allowed to access the host ports/protocols"`
}

type PortSelectorWrapper struct {
	Ports    []any  `yaml:"ports" jsonschema:"description=List of ports or port ranges"`
	Protocol string `yaml:"protocol" jsonschema:"description=Protocol (can be tcp or udp)"`
}

type IngressConfigWrapper struct {
	Subnet string `yaml:"subnet" jsonschema:"description=Source subnet"`
	Except string `yaml:"except" jsonschema:"description=Source subnet to exclude from the subnet"`
}

func (Node) JSONSchemaProperty(prop string) any {
	if prop == "installDiskSelector" {
		return &InstallDiskSelectorWrapper{}
	}
	return nil
}

func (IngressFirewall) JSONSchemaAlias() any {
	return &IngressFirewallWrapper{}
}
