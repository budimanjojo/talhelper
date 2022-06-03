package config

import (
	"encoding/json"

	talosconfig "github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
	"github.com/talos-systems/talos/pkg/machinery/constants"
)

func ParseTalosInput(config TalhelperConfig) (*generate.Input, error) {
	kubernetesVersion := config.KubernetesVersion

	if kubernetesVersion == "" {
		kubernetesVersion = constants.DefaultKubernetesVersion
	}

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

	var patch []map[string]interface{}
	var iPatch map[string]interface{}

	switch node.ControlPlane {
	case true:
		cfg, err = generate.Config(machine.TypeControlPlane, input)
		if err != nil {
			return nil, err
		}
		patch = config.ControlPlane.ConfigPatches
		iPatch = config.ControlPlane.InlinePatch
	case false:
		cfg, err = generate.Config(machine.TypeWorker, input)
		if err != nil {
			return nil, err
		}
		patch = config.Worker.ConfigPatches
		iPatch = config.Worker.InlinePatch
	}

	cfg.MachineConfig.MachineInstall.InstallDisk = node.InstallDisk

	if node.Domain == "" {
		cfg.MachineConfig.MachineNetwork.NetworkHostname = node.Hostname
	} else {
		cfg.MachineConfig.MachineNetwork.NetworkHostname = node.Hostname + "." + node.Domain
	}

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

	if iPatch != nil {
		marshaledIPatch, err := json.Marshal(iPatch)
		if err != nil {
			return nil, err
		}

		finalCfg, err := ApplyInlinePatchFromYaml(marshaledIPatch, patchedCfg)
		if err != nil {
			return nil, err
		}

		return finalCfg, nil
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
