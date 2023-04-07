package talos

import (
	"github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
)

// LoadTalosConfigFromFile generates Talos `config.Provider` from a file path.
// It also returns an error, if any.
func LoadTalosConfigFromFile(cfgFile string) (config.Provider, error) {
	return configloader.NewFromFile(cfgFile)
}

// LoadTalosConfig generates Talos `config.Provider` from bytes.
// It also returns an error, if any.
func LoadTalosConfig(cfgFile []byte) (config.Provider, error) {
	return configloader.NewFromBytes(cfgFile)
}

// IsControlPlane returns true if `c` is a controlplane.
func IsControlPlane(c config.Provider) bool {
	return c.Machine().Type().String() == "controlplane"
}
