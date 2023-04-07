package generate

import (
	"github.com/budimanjojo/talhelper/pkg/secret"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1/generate"
)

// GenerateSecret generates `SecretsBundle` in the specified path.
// It returns an error, if any.
func GenerateSecret(cfg string) error {
	var s *generate.SecretsBundle
	var err error
	switch cfg {
	case "":
		s, err = talos.NewSecretBundle(generate.NewClock())
		if err != nil {
			return err
		}
	default:
		cfg, err := talos.LoadTalosConfigFromFile(cfg)
		if err != nil {
			return err
		}
		s = talos.NewSecretBundleFromCfg(generate.NewClock(), cfg)
	}

	err = secret.PrintSecretBundle(s)
	if err != nil {
		return err
	}

	return nil
}
