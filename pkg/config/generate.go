package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	talosconfig "github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
	"github.com/talos-systems/talos/pkg/machinery/constants"
)

func (config TalhelperConfig) GenerateConfig(outputDir string) error {
	const (
		kubernetesVersion = constants.DefaultKubernetesVersion
	)

	versionContract, err := talosconfig.ParseContractFromVersion(config.TalosVersion)
	if err != nil {
		return fmt.Errorf("failed to parse version contract: %s", err)
	}

	secrets, err := generate.NewSecretsBundle(generate.NewClock(), generate.WithVersionContract(versionContract))
	if err != nil {
		return fmt.Errorf("failed to generate secrets bundle: %s", err)
	}

	input, err := generate.NewInput(config.ClusterName, config.Endpoint, kubernetesVersion, secrets, generate.WithVersionContract(versionContract))
	if err != nil {
		return fmt.Errorf("failed to generate input: %s", err)
	}

	for _, node := range config.Nodes {
		var cfg *v1alpha1.Config
		var patch []byte

		switch node.ControlPlane {
		case true:
			cfg, err = generate.Config(machine.TypeControlPlane, input)
			if err != nil {
				return fmt.Errorf("failed to generate config for node %q: %s", node.Hostname, err)
			}
			patch, err = json.Marshal(config.ControlPlane.ConfigPatches)
			if err != nil {
				return fmt.Errorf("failed to decode patch for node %q: %s", node.Hostname, err)
			}
		case false:
			cfg, err = generate.Config(machine.TypeWorker, input)
			if err != nil {
				return fmt.Errorf("failed to generate config for node %q: %s", node.Hostname, err)
			}
			patch, err = json.Marshal(config.ControlPlane.ConfigPatches)
			if err != nil {
				return fmt.Errorf("failed to decode patch for node %q: %s", node.Hostname, err)
			}
		}

		cfg.MachineConfig.MachineInstall.InstallDisk = node.InstallDisk
		cfg.MachineConfig.MachineNetwork.NetworkHostname = node.Hostname

		marshaledCfg, err := cfg.Bytes()
		if err != nil {
			return fmt.Errorf("failed to generate config for node %q: %s", node.Hostname, err)
		}

		cfgFile := outputDir + "/" + config.ClusterName + node.Hostname + ".yaml"

		patchedCfgFile, err := applyPatchFromYaml(patch, marshaledCfg)
		if err != nil {
			return fmt.Errorf("failed to apply patch for node %q: %s", node.Hostname, err)
		}

		err = dumpConfig(cfgFile, patchedCfgFile)
		if err != nil {
			return fmt.Errorf("failed to dump config for node %q: %s", node.Hostname, err)
		}

		fmt.Printf("generated config for %s in %s\n", node.Hostname, cfgFile)
	}

	var endpointList []string
	for _, node := range config.Nodes {
		endpointList = append(endpointList, node.IPAddress)
	}

	clientCfg, err := generate.Talosconfig(input, generate.WithEndpointList(endpointList))
	if err != nil {
		return fmt.Errorf("failed to generate client config: %s", err)
	}

	marshaledClientCfg, err := clientCfg.Bytes()
	if err != nil {
		return fmt.Errorf("failed to generate client config: %s", err)
	}

	err = dumpConfig(outputDir+"/talosconfig", marshaledClientCfg)
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
