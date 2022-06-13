package talos

import (
	"github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/configloader"
)

func LoadTalosConfig(cfgFile []byte) (config.Provider, error) {
	return configloader.NewFromBytes(cfgFile)
}

func IsControlPlane(c config.Provider) bool {
	return c.Machine().Type().String() == "controlplane"
}
