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
	input, err := parseTalosInput(config)
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

		err = createGitIgnore(outputDir, fileName)
		if err != nil {
			return fmt.Errorf("failed to create gitignore file for node %q: %s", node.Hostname, err)
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

	err = createGitIgnore(outputDir, fileName)
	if err != nil {
		return fmt.Errorf("failed to create client config gitignore file: %s", err)
	}

	fmt.Printf("generated client config in %s\n", outputDir+"/talosconfig")

	return nil
}

func parseTalosInput(config TalhelperConfig) (*generate.Input, error) {
	const (
		kubernetesVersion = constants.DefaultKubernetesVersion
	)

	versionContract, err := talosconfig.ParseContractFromVersion(config.TalosVersion)
	if err != nil {
		return nil, err
	}

	secrets, err := generate.NewSecretsBundle(generate.NewClock(), generate.WithVersionContract(versionContract))
	if err != nil {
		return nil, err
	}

	input, err := generate.NewInput(config.ClusterName, config.Endpoint, kubernetesVersion, secrets, generate.WithVersionContract(versionContract))
	if err != nil {
		return nil, err
	}

	return input, nil
}

func createTalosClusterConfig(node nodes, config TalhelperConfig, input *generate.Input) (CfgFile []byte, err error) {
	var cfg *v1alpha1.Config

	controlPlanePatch := mergePatchSlices(config.ControlPlane.Patches, config.ControlPlane.EncryptedPatches)
	workerPatch := mergePatchSlices(config.Worker.Patches, config.Worker.EncryptedPatches)
	var patch []map[string]interface{}

	switch node.ControlPlane {
	case true:
		cfg, err = generate.Config(machine.TypeControlPlane, input)
		if err != nil {
			return nil, err
		}
		patch = controlPlanePatch
	case false:
		cfg, err = generate.Config(machine.TypeWorker, input)
		if err != nil {
			return nil, err
		}
		patch = workerPatch
	}

	cfg.MachineConfig.MachineInstall.InstallDisk = node.InstallDisk
	cfg.MachineConfig.MachineNetwork.NetworkHostname = node.Hostname

	marshaledCfg, err := cfg.Bytes()
	if err != nil {
		return nil, err
	}

	marshaledPatch, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	patchedCfg, err := applyPatchFromYaml(marshaledPatch, marshaledCfg)
	if err != nil {
		return nil, err
	}
	return patchedCfg, nil
}

func createTalosClientConfig(config TalhelperConfig, input *generate.Input) ([]byte, error) {
	var endpointList []string
	for _, node := range config.Nodes {
		endpointList = append(endpointList, node.IPAddress)
	}

	clientCfg, err := generate.Talosconfig(input, generate.WithEndpointList(endpointList))
	if err != nil {
		return nil, err
	}

	marshaledClientCfg, err := clientCfg.Bytes()
	if err != nil {
		return nil, err
	}

	return marshaledClientCfg, nil
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
