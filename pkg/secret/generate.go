package secret

import (
	"bytes"
	"os"

	talconfig "github.com/budimanjojo/talhelper/pkg/config"
	"github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"gopkg.in/yaml.v3"
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
	err = yaml.Unmarshal(yamlCfg, &m)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	err = encoder.Encode(m)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, buf.Bytes(), 0700)
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
