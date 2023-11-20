package parse

import (
	"fmt"
	"os"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/substitute"
	"github.com/fatih/color"
)

func ParseConfig(configFile string, envFile []string) (config.TalhelperConfig, error) {
	cfgByte, err := os.ReadFile(configFile)
	if err != nil {
		return config.TalhelperConfig{}, fmt.Errorf("failed to read config file: %s", err)
	}

	if err := substitute.LoadEnvFromFiles(envFile); err != nil {
		return config.TalhelperConfig{}, fmt.Errorf("failed to load env file: %s", err)
	}

	cfgByte, err = substitute.SubstituteEnvFromByte(cfgByte)
	if err != nil {
		return config.TalhelperConfig{}, fmt.Errorf("failed to substitute env: %s", err)
	}

	cfg, err := config.NewFromByte(cfgByte)
	if err != nil {
		return config.TalhelperConfig{}, fmt.Errorf("failed to unmarshal config file: %s", err)
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
			os.Exit(1)
		}
	}

	return cfg, nil
}