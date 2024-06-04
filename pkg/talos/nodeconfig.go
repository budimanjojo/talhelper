package talos

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/image-factory/pkg/schematic"
	taloscfg "github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/machine"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
)

func GenerateNodeConfigBytes(node *config.Node, input *generate.Input, iFactory *config.ImageFactory, offlineMode bool) ([]byte, error) {
	cfg, err := GenerateNodeConfig(node, input, iFactory, offlineMode)
	if err != nil {
		return nil, err
	}
	return cfg.Bytes()
}

func GenerateNodeConfig(node *config.Node, input *generate.Input, iFactory *config.ImageFactory, offlineMode bool) (taloscfg.Provider, error) {
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
		slog.Debug("encryption with secretbox is supported and AESCBCEncryptionSecret is not empty")
		c.RawV1Alpha1().ClusterConfig.ClusterAESCBCEncryptionSecret = input.Options.SecretsBundle.Secrets.AESCBCEncryptionSecret
	}

	cfg := applyNodeOverride(node, c)

	installerURL, err := installerURL(node, c, iFactory, offlineMode)
	if err != nil {
		return nil, err
	}
	slog.Debug(fmt.Sprintf("installer URL for %s is set to: %s", node.Hostname, installerURL))
	cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallImage = installerURL

	return cfg, nil
}

func applyNodeOverride(node *config.Node, cfg taloscfg.Provider) taloscfg.Provider {
	if !node.IgnoreHostname {
		slog.Debug(fmt.Sprintf("setting hostname to %s", node.Hostname))
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkHostname = node.Hostname
	}

	if len(node.Nameservers) > 0 {
		slog.Debug(fmt.Sprintf("setting nameservers to %s", node.Nameservers))
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NameServers = node.Nameservers
	}

	if node.DisableSearchDomain {
		slog.Debug("setting disableSearchDomain to true")
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkDisableSearchDomain = &node.DisableSearchDomain
	}

	if len(node.NetworkInterfaces) > 0 {
		slog.Debug("setting network interfaces")
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkInterfaces = node.NetworkInterfaces
	}

	if node.InstallDisk != "" {
		slog.Debug(fmt.Sprintf("setting install disk to %s", node.InstallDisk))
		cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallDisk = node.InstallDisk
	}

	if node.InstallDiskSelector != nil {
		slog.Debug("setting install disk selector")
		cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallDiskSelector = node.InstallDiskSelector
	}

	if len(node.MachineDisks) > 0 {
		slog.Debug("setting machine disks")
		cfg.RawV1Alpha1().MachineConfig.MachineDisks = node.MachineDisks
	}

	if len(node.KernelModules) > 0 {
		slog.Debug("setting kernel modules")
		cfg.RawV1Alpha1().MachineConfig.MachineKernel = &v1alpha1.KernelConfig{}
		cfg.RawV1Alpha1().MachineConfig.MachineKernel.KernelModules = node.KernelModules
	}

	if node.NodeLabels != nil {
		slog.Debug(fmt.Sprintf("setting node labels to %s", node.NodeLabels))
		cfg.RawV1Alpha1().MachineConfig.MachineNodeLabels = node.NodeLabels
	}

	if node.NodeTaints != nil {
		slog.Debug(fmt.Sprintf("setting node taints to %s", node.NodeTaints))
		cfg.RawV1Alpha1().MachineConfig.MachineNodeTaints = node.NodeTaints
	}

	if len(node.MachineFiles) > 0 {
		slog.Debug("setting machine files")
		cfg.RawV1Alpha1().MachineConfig.MachineFiles = node.MachineFiles
	}

	if node.Schematic != nil && len(node.Schematic.Customization.ExtraKernelArgs) > 0 {
		slog.Debug("appending schematic kernel args to install kernel args")
		cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallExtraKernelArgs = append(cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallExtraKernelArgs, node.Schematic.Customization.ExtraKernelArgs...)
	}

	return cfg
}

func installerURL(node *config.Node, cfg taloscfg.Provider, iFactory *config.ImageFactory, offlineMode bool) (string, error) {
	version := strings.Split(cfg.Machine().Install().Image(), ":")

	if node.Schematic != nil && node.TalosImageURL == "" {
		url, err := GetInstallerURL(node.Schematic, iFactory, node.GetMachineSpec(), version[1], offlineMode)
		if err != nil {
			return "", err
		}
		return url, nil
	}

	if node.TalosImageURL != "" {
		return node.TalosImageURL + ":" + version[1], nil
	}

	return GetInstallerURL(&schematic.Schematic{}, iFactory, node.GetMachineSpec(), version[1], offlineMode)
}
