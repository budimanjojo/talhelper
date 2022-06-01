package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func (config TalhelperConfig) GenerateConfig(outputDir string) error {
	input, err := ParseTalosInput(config)
	if err != nil {
		return fmt.Errorf("failed to generate talos input: %s", err)
	}

	for _, node := range config.Nodes {
		fileName := config.ClusterName + "-" + node.Hostname + ".yaml"
		cfgFile := outputDir + "/" + fileName

		patchedCfg, err := createTalosClusterConfig(node, config, input)
		if err != nil {
			return fmt.Errorf("failed to create Talos cluster config: %s", err)
		}

		err = dumpConfig(cfgFile, patchedCfg)
		if err != nil {
			return fmt.Errorf("failed to dump config for node %q: %s", node.Hostname, err)
		}

		fmt.Printf("generated config for %s in %s\n", node.Hostname, cfgFile)
	}

	marshaledClientCfg, err := createTalosClientConfig(config, input)
	if err != nil {
		return fmt.Errorf("failed to create Talos client config: %s", err)
	}

	fileName := "talosconfig"
	err = dumpConfig(outputDir+"/"+fileName, marshaledClientCfg)
	if err != nil {
		return fmt.Errorf("failed to dump client config: %s", err)
	}

	fmt.Printf("generated client config in %s\n", outputDir+"/talosconfig")

	return nil
}

func dumpConfig(filePath string, marshaledCfg []byte) error {
	dirName := filepath.Dir(filePath)

	_, err := os.Stat(dirName)
	if err != nil {
		err := os.MkdirAll(dirName, 0700)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(filePath, marshaledCfg, 0600)
	if err != nil {
		return err
	}

	return nil
}
