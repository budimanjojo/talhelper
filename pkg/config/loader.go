package config

import (
	"fmt"
	"os"

	"github.com/budimanjojo/talhelper/pkg/substitute"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

// LoadAndValidateFromFile takes a file path and yaml encoded env files path, do envsubst
// from envPaths. The resulted TalhelperConfig will be validated before being returned.
// It returns an error, if any.
func LoadAndValidateFromFile(filePath string, envPaths []string) (*TalhelperConfig, error) {
	cfgByte, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	if err := substitute.LoadEnvFromFiles(envPaths); err != nil {
		return nil, fmt.Errorf("failed to load env file: %s", err)
	}

	cfgByte, err = substitute.SubstituteEnvFromByte(cfgByte)
	if err != nil {
		return nil, fmt.Errorf("failed to substitute env: %s", err)
	}

	cfg, err := NewFromByte(cfgByte)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %s", err)
	}

	for k, node := range cfg.Nodes {
		switch node.ControlPlane {
		case true:
			cfg.Nodes[k].OverrideGlobalCfg(cfg.ControlPlane)
		case false:
			cfg.Nodes[k].OverrideGlobalCfg(cfg.Worker)
		}
	}

	errs, warns := cfg.Validate()
	if len(errs) > 0 || len(warns) > 0 {
		color.Red("There are issues with your talhelper config file:")
		grouped := make(map[string][]string)
		for _, v := range errs {
			grouped[v.Field] = append(grouped[v.Field], v.Message.Error())
		}
		for _, v := range warns {
			grouped[v.Field] = append(grouped[v.Field], v.Message)
		}
		for field, list := range grouped {
			color.Yellow("field: %q\n", field)
			for _, l := range list {
				fmt.Printf(l + "\n")
			}
		}

		if len(errs) > 0 {
			return nil, fmt.Errorf("please fix issues with your config file")
		}
	}

	return cfg, nil
}

// NewFromByte takes bytes and convert it into Talhelper config.
// It also returns an error, if any.
func NewFromByte(source []byte) (*TalhelperConfig, error) {
	return newConfig(source)
}

// NewFromFile takes a file path and convert the contents into Talhelper config.
// It also returns an error, if any.
func NewFromFile(path string) (c *TalhelperConfig, err error) {
	source, err := fromFile(path)
	if err != nil {
		return c, err
	}

	return newConfig(source)
}

// fromFile is a wrapper for `os.ReadFile`.
func fromFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// newConfig takes bytes and convert it into Talhelper config.
// It also returns an error, if any.
func newConfig(source []byte) (c *TalhelperConfig, err error) {
	err = yaml.Unmarshal(source, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
