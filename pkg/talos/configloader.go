package talos

import (
	"github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
)

func LoadTalosConfigFromFile(cfgFile string) (config.Provider, error) {
	return configloader.NewFromFile(cfgFile)
}

func LoadTalosConfig(cfgFile []byte) (config.Provider, error) {
	return configloader.NewFromBytes(cfgFile)
}

func IsControlPlane(c config.Provider) bool {
	return c.Machine().Type().String() == "controlplane"
}
