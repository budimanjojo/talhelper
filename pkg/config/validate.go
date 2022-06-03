package config

import (
	"fmt"

	"github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/configloader"
)

type mode int

func (m mode) String() string {
	return ""
}

func (m mode) RequiresInstall() bool {
	return m == 2
}

func parseMode(s string) (mod mode, err error) {
	switch s {
	case "cloud":
		mod = 0
	case "container":
		mod = 1
	case "metal":
		mod = 2
	default:
		return mod, fmt.Errorf("unknown Talos runtime mode: %q", s)
	}

	return mod, nil
}

func validateConfig(cfgFile []byte, mode string) error {
	cfg, err := configloader.NewFromBytes(cfgFile)
	if err != nil {
		return err
	}

	m, err := parseMode(mode)
	if err != nil {
		return err
	}

	opts := []config.ValidationOption{config.WithLocal()}

	warnings, err := cfg.Validate(m, opts...)
	for _, w := range warnings {
		fmt.Printf("%s\n", w)
	}
	if err != nil {
		return err
	}

	return nil
}
