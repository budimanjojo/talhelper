package generate

import (
	"os"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/secret"
	"github.com/budimanjojo/talhelper/pkg/talos"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"gopkg.in/yaml.v3"
)

func GenerateOutput(cfg string) error {
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
		s = talos.NewSecretFromCfg(generate.NewClock(), cfg)
	}

	secret.PrintSortedSecrets(s)
	return nil
}

func PatchTalhelperConfig(cfgFile string) error {
	cfg, err := os.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(cfg, &m); err != nil {
		return err
	}

	cfg, err = m.ApplyInlinePatch([]byte(secret.SecretPatch))
	if err != nil {
		return err
	}

	cfg, err = m.Encode(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(cfgFile, cfg, 0700)
	if err != nil {
		return err
	}

	return nil
}
