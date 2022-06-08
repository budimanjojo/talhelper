package config

import (
	"fmt"
	"os"
	"path/filepath"

	talosconfig "github.com/talos-systems/talos/pkg/machinery/config"
	"sigs.k8s.io/yaml"
)

func (config TalhelperConfig) GenerateConfig(outputDir, mode string) error {
	var cfgDump talosconfig.Provider
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

		if node.InlinePatch != nil {
			iPatch, err := yaml.Marshal(node.InlinePatch)
			if err != nil {
				return fmt.Errorf("failed to marshal node inline patch for node %q: %s", node.Hostname, err)
			}

			patchedCfg, err = ApplyInlinePatchFromYaml(iPatch, patchedCfg)
			if err != nil {
				return fmt.Errorf("failed to apply node patch for node %q: %s", node.Hostname, err)
			}
		}

		if len(node.ConfigPatches) != 0 {
			nodePatches, err := yaml.Marshal(node.ConfigPatches)
			if err != nil {
				return fmt.Errorf("failed to marshal node configPatches for node %q: %s", node.Hostname, err)
			}

			patchedCfg, err = applyPatchFromYaml(nodePatches, patchedCfg)
			if err != nil {
				return fmt.Errorf("failed to apply node configPatches for node %q: %s", node.Hostname, err)
			}
		}

		cfgDump, err = parseTalosConfig(patchedCfg)
		if err != nil {
			return fmt.Errorf("failed to dump config for node %q: %s", node.Hostname, err)
		}

		err = validateConfig(patchedCfg, mode)
		if err != nil {
			return fmt.Errorf("failed to verify config for node %q: %s", node.Hostname, err)
		}

		err = dumpConfig(cfgFile, patchedCfg)
		if err != nil {
			return fmt.Errorf("failed to dump config for node %q: %s", node.Hostname, err)
		}

		fmt.Printf("generated config for %s in %s\n", node.Hostname, cfgFile)
	}

	machineCert := cfgDump.Machine().Security().CA()

	marshaledClientCfg, err := createTalosClientConfig(config, input, machineCert)
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
