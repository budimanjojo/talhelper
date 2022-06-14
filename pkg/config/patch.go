package config

import (
	"github.com/budimanjojo/talhelper/pkg/patcher"
	"gopkg.in/yaml.v3"
)

func (c *TalhelperConfig) ApplyInlinePatch(patch []byte) ([]byte, error) {
	cfg, err := yaml.Marshal(c)
	if err != nil {
		return nil, err
	}

	cfg, err = patcher.JSON7396FromYAML(patch, cfg)
	if err != nil {
		return nil, err
	}

	var m TalhelperConfig
	cfg, err = m.Encode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
