package talos

import (
	"fmt"
	"os"

	"github.com/siderolabs/talos/pkg/machinery/config/validation"
)

type mode int

const (
	modeCloud mode = iota
	modeContainer
	modeMetal
)

func (m mode) String() string {
	return ""
}

func (m mode) RequiresInstall() bool {
	return m == modeMetal
}

func (m mode) InContainer() bool {
	return m == modeContainer
}

// parseMode takes string and convert it to `mode`.
// It also returns an error, if any.
func parseMode(s string) (mod mode, err error) {
	switch s {
	case "cloud":
		mod = modeCloud
	case "container":
		mod = modeContainer
	case "metal":
		mod = modeMetal
	default:
		return mod, fmt.Errorf("unknown Talos runtime mode: %q", s)
	}

	return mod, nil
}

// ValidateConfigFromFile returns an error if file path is not a valid
// Talos configuration for the specified `mode`.
func ValidateConfigFromFile(path, mode string) error {
	output, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return ValidateConfigFromBytes(output, mode)
}

// ValidateConfigFromBytes returns an error if `cfgFile` is not a valid
// Talos configuration for the specified `mode`.
func ValidateConfigFromBytes(cfgFile []byte, mode string) error {
	cfg, err := LoadTalosConfig(cfgFile)
	if err != nil {
		return err
	}

	m, err := parseMode(mode)
	if err != nil {
		return err
	}

	warnings, err := cfg.Validate(m, validation.WithLocal(), validation.WithStrict())
	for _, w := range warnings {
		fmt.Printf("%s\n", w)
	}
	if err != nil {
		return err
	}

	return nil
}
