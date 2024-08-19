package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/budimanjojo/talhelper/v3/pkg/substitute"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

// LoadAndValidateFromFile takes a file path and yaml encoded env files path, do envsubst
// from envPaths. The resulted TalhelperConfig will be validated before being returned.
// It returns an error, if any.
func LoadAndValidateFromFile(filePath string, envPaths []string, showWarns bool) (*TalhelperConfig, error) {
	slog.Debug("start loading and validating config file")
	slog.Debug(fmt.Sprintf("reading %s", filePath))
	cfgByte, err := FromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	if err := substitute.LoadEnvFromFiles(envPaths); err != nil {
		return nil, fmt.Errorf("failed to load env file: %s", err)
	}

	slog.Debug("substituting config file with environment variable")
	cfgByte, err = substitute.SubstituteEnvFromByte(cfgByte)
	if err != nil {
		return nil, fmt.Errorf("failed to substitute env: %s", err)
	}

	cfg, err := NewFromByte(cfgByte)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %s", err)
	}

	for k := range cfg.Nodes {
		node := &cfg.Nodes[k]

		switch node.ControlPlane {
		case true:
			slog.Debug(fmt.Sprintf("overriding global controlplane node config for %s", node.Hostname))
			node.OverrideGlobalCfg(cfg.ControlPlane)
		case false:
			slog.Debug(fmt.Sprintf("overriding global worker node config for %s", node.Hostname))
			node.OverrideGlobalCfg(cfg.Worker)
		}

		if len(node.MachineFiles) > 0 {
			for i, file := range node.MachineFiles {
				contents, err := ensureFileContent(file.FileContent)
				if err != nil {
					return nil, fmt.Errorf("failed to get machine file content for %s in `machineFiles[%d]`: %s", node.Hostname, i, err)
				}
				file.FileContent = contents
			}
		}
	}

	errs, warns := cfg.Validate()
	grouped := make(map[string][]string)

	for _, v := range errs {
		grouped[v.Field] = append(grouped[v.Field], v.Message.Error())
	}

	if showWarns {
		for _, v := range warns {
			grouped[v.Field] = append(grouped[v.Field], v.Message)
		}
	}

	if len(grouped) > 0 {
		color.Red("There are issues with your Talhelper config file:")

		for field, list := range grouped {
			color.Yellow("field: %q\n", field)
			for _, l := range list {
				fmt.Println(l)
			}
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("please fix issues with your config file")
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
	source, err := FromFile(path)
	if err != nil {
		return c, err
	}

	return newConfig(source)
}

// FromFile is a wrapper for `os.ReadFile` with modified error if path doesn't exist.
func FromFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("%s doesn't exist. Refer to this docs for more information on how to create one: https://budimanjojo.github.io/talhelper/latest/guides/#example-talconfigyaml", path)
	}
	return b, err
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

func ensureFileContent(value string) (string, error) {
	if strings.HasPrefix(value, "@") {
		slog.Debug(fmt.Sprintf("getting file content of %s", value))
		filename := value[1:]

		contents, err := os.ReadFile(filename)
		if err != nil {
			return "", err
		}

		substituted, err := substitute.SubstituteEnvFromByte(contents)
		if err != nil {
			return "", err
		}

		return string(substituted), nil
	}

	return value, nil
}
