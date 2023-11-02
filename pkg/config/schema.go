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

func (Node) JSONSchemaProperty(prop string) any {
	if prop == "installDiskSelector" {
		return &InstallDiskSelectorWrapper{}
	}
	return nil
}
