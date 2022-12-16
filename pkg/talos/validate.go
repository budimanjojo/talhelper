package talos

import (
	"fmt"
	"io/ioutil"

	"github.com/siderolabs/talos/pkg/machinery/config"
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

func ValidateConfigFromFile(path, mode string) error {
	output, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return ValidateConfigFromBytes(output, mode)
}

func ValidateConfigFromBytes(cfgFile []byte, mode string) error {
	cfg, err := LoadTalosConfig(cfgFile)
	if err != nil {
		return err
	}

	m, err := parseMode(mode)
	if err != nil {
		return err
	}

	opts := []config.ValidationOption{config.WithLocal(), config.WithStrict()}

	warnings, err := cfg.Validate(m, opts...)
	for _, w := range warnings {
		fmt.Printf("%s\n", w)
	}
	if err != nil {
		return err
	}

	return nil
}
