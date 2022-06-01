package secret

import (
	"os"

	talconfig "github.com/budimanjojo/talhelper/pkg/config"
)

func GenerateSecret(config talconfig.TalhelperConfig, configFile string) error {
	cf, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	yamlCfg, err := talconfig.ApplyInlinePatchFromYaml([]byte(secretPatch), cf)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, yamlCfg, 0700)
	if err != nil {
		return err
	}

	return nil
}
