package generate

import (
	talhelperCfg "github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/secret"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate/secrets"
)

// GenerateSecret generates `SecretsBundle` in the specified path.
// It returns an error, if any.
func GenerateSecret(cfg string) error {
	var s *secrets.Bundle
	var err error
	switch cfg {
	case "":
		version, _ := config.ParseContractFromVersion(talhelperCfg.LatestTalosVersion)
		s, err = talos.NewSecretBundle(secrets.NewClock(), *version)
		if err != nil {
			return err
		}
	default:
		cfg, err := talos.LoadTalosConfigFromFile(cfg)
		if err != nil {
			return err
		}
		s = talos.NewSecretBundleFromCfg(secrets.NewClock(), cfg)
	}

	err = secret.PrintSecretBundle(s)
	if err != nil {
		return err
	}

	return nil
}
