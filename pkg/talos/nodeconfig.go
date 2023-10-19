package talos

import (
	"strings"

	"github.com/budimanjojo/talhelper/pkg/config"
	taloscfg "github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/machine"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
)

func GenerateNodeConfigBytes(node *config.Node, input *generate.Input) ([]byte, error) {
	cfg, err := GenerateNodeConfig(node, input)
	if err != nil {
		return nil, err
	}
	return cfg.Bytes()
}

func GenerateNodeConfig(node *config.Node, input *generate.Input) (taloscfg.Provider, error) {
	var c taloscfg.Provider
	var err error

	switch node.ControlPlane {
	case true:
		c, err = input.Config(machine.TypeControlPlane)
		if err != nil {
			return nil, err
		}
	case false:
		c, err = input.Config(machine.TypeWorker)
		if err != nil {
			return nil, err
		}
	}

	// https://github.com/budimanjojo/talhelper/issues/81
	if input.Options.VersionContract.SecretboxEncryptionSupported() && input.Options.SecretsBundle.Secrets.AESCBCEncryptionSecret != "" {
		c.RawV1Alpha1().ClusterConfig.ClusterAESCBCEncryptionSecret = input.Options.SecretsBundle.Secrets.AESCBCEncryptionSecret
	}

	cfg := applyNodeOverride(node, c)

	return cfg, nil
}

func applyNodeOverride(node *config.Node, cfg taloscfg.Provider) taloscfg.Provider {
	cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkHostname = node.Hostname

	if node.TalosImageURL != "" {
		version := strings.Split(cfg.Machine().Install().Image(), ":")
		cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallImage = node.TalosImageURL + ":" + version[1]
	}

	if len(node.Nameservers) > 0 {
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NameServers = node.Nameservers
	}

	if node.DisableSearchDomain {
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkDisableSearchDomain = &node.DisableSearchDomain
	}

	if len(node.NetworkInterfaces) > 0 {
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkInterfaces = node.NetworkInterfaces
	}

	if node.InstallDisk != "" {
		cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallDisk = node.InstallDisk
	}

	if node.InstallDiskSelector != nil {
		cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallDiskSelector = node.InstallDiskSelector
	}

	if len(node.MachineDisks) > 0 {
		cfg.RawV1Alpha1().MachineConfig.MachineDisks = node.MachineDisks
	}

	if len(node.KernelModules) > 0 {
		cfg.RawV1Alpha1().MachineConfig.MachineKernel = &v1alpha1.KernelConfig{}
		cfg.RawV1Alpha1().MachineConfig.MachineKernel.KernelModules = node.KernelModules
	}

	if node.NodeLabels != nil {
		cfg.RawV1Alpha1().MachineConfig.MachineNodeLabels = node.NodeLabels
	}

	if len(node.Extensions) > 0 {
		cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallExtensions = node.Extensions
	}

	if len(node.MachineFiles) > 0 {
		cfg.RawV1Alpha1().MachineConfig.MachineFiles = node.MachineFiles
	}

	return cfg
}
