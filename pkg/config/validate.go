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
	return true
}

func parseMode() (mode, error) {
	return 1, nil
}

func ValidateConfig(cfgFile string) error {
	cfg, err := configloader.NewFromFile(cfgFile)
	if err != nil {
		return err
	}

	mode, err := parseMode()

	opts := []config.ValidationOption{config.WithLocal()}

	warnings, err := cfg.Validate(mode, opts...)
	for _, w := range warnings {
		fmt.Printf("%s\n", w)
	}
	if err != nil {
		return err
	}

	return nil
}
