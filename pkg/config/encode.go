package config

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func (c *TalhelperConfig) Encode(cfg []byte) ([]byte, error) {
	err := yaml.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	err = encoder.Encode(c)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
