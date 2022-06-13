package secret

import (
	"os"

	talconfig "github.com/budimanjojo/talhelper/pkg/config"
	"github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
)

func PatchTalconfig(configFile string) error {
	cf, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	yamlCfg, err := talconfig.ApplyInlinePatchFromYaml([]byte(secretPatch), cf)
	if err != nil {
		return err
	}

	// Reencode so the formatting stays
	var m talconfig.TalhelperConfig
	cfg, err := m.Encode(yamlCfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, cfg, 0700)
	if err != nil {
		return err
	}

	return nil
}

func NewSecretFromCfg(clock generate.Clock, talosCfg config.Provider) *generate.SecretsBundle {
	return generate.NewSecretsBundleFromConfig(clock, talosCfg)
}

func NewSecretBundle(clock generate.Clock) (*generate.SecretsBundle, error) {
	return generate.NewSecretsBundle(clock)
}
