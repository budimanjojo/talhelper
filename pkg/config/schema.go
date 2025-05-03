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

type DiskSelectorWrapper struct {
	Match string `yaml:"match" jsonschema:"description=The Common Expression Language (CEL) expression to match the disk"`
}

type ProvisioningSpecWrapper struct {
	DiskSelectorSpec    DiskSelectorWrapper `yaml:"diskSelector" jsonschema:"description=The disk selector expression"`
	ProvisioningGrow    bool                `yaml:"grow" jsonschema:"description=Should the volume grow to the size of the disk (if possible)"`
	ProvisioningMinSize string              `yaml:"minSize" jsonschema:"description=The minimum size of the volume,example=2.5GiB"`
	ProvisioningMaxSize string              `yaml:"maxSize" jsonschema:"description=The maximum size of the volume, if not specified the volume can grow to the size of the disk,example=50GiB"`
}

type FilesystemSpecWrapper struct {
	FilesystemType string `yaml:"type" jsonschema:"default=xfs,description=Filesystem type,enum=ext4,enum=xfs"`
}

func (Node) JSONSchemaProperty(prop string) any {
	if prop == "installDiskSelector" {
		return &InstallDiskSelectorWrapper{}
	}
	return nil
}

func (UserVolume) JSONSchemaProperty(prop string) any {
	if prop == "provisioning" {
		return &ProvisioningSpecWrapper{}
	}
	if prop == "filesystem" {
		return &FilesystemSpecWrapper{}
	}
	return nil
}

func (Volume) JSONSchemaProperty(prop string) any {
	if prop == "provisioning" {
		return &ProvisioningSpecWrapper{}
	}
	return nil
}

func (IngressFirewall) JSONSchemaAlias() any {
	return &IngressFirewallWrapper{}
}
