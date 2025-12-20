package talos

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/templating"
	"github.com/siderolabs/image-factory/pkg/schematic"
	taloscfg "github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/encoder"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/machine"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
)

func GenerateNodeConfigBytes(node *config.Node, input *generate.Input, iFactory *config.ImageFactory, offlineMode bool) ([]byte, error) {
	cfg, err := GenerateNodeConfig(node, input, iFactory, offlineMode)
	if err != nil {
		return nil, err
	}
	return cfg.EncodeBytes(encoder.WithComments(encoder.CommentsDisabled))
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

	if !input.Options.VersionContract.MultidocNetworkConfigSupported() && !node.IgnoreHostname {
		slog.Debug(fmt.Sprintf("setting hostname to %s", node.Hostname))
		//nolint:staticcheck
		c.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkHostname = node.Hostname
	}

	if !input.Options.VersionContract.MultidocNetworkConfigSupported() && len(node.Nameservers) > 0 {
		slog.Debug(fmt.Sprintf("setting nameservers to %s", node.Nameservers))
		//nolint:staticcheck
		c.RawV1Alpha1().MachineConfig.MachineNetwork.NameServers = node.Nameservers
	}

	if !input.Options.VersionContract.MultidocNetworkConfigSupported() && node.DisableSearchDomain {
		slog.Debug("setting disableSearchDomain to true")
		//nolint:staticcheck
		c.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkDisableSearchDomain = &node.DisableSearchDomain
	}

	cfg := applyNodeOverride(node, c, *input.Options.VersionContract)

	installerURL, err := installerURL(node, c, iFactory, offlineMode)
	if err != nil {
		return nil, err
	}
	slog.Debug(fmt.Sprintf("installer URL for %s is set to: %s", node.Hostname, installerURL))
	cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallImage = installerURL

	// Templating should be done as late as possible to maximize the amount of infomation
	// available for templating
	cfg, err = templateConfig(node, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func applyNodeOverride(node *config.Node, cfg taloscfg.Provider, vc taloscfg.VersionContract) taloscfg.Provider {
	if len(node.NetworkInterfaces) > 0 {
		slog.Debug("setting network interfaces")
		//nolint:staticcheck
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
		//nolint:all
		cfg.RawV1Alpha1().MachineConfig.MachineDisks = node.MachineDisks
	}

	if len(node.KernelModules) > 0 {
		slog.Debug("setting kernel modules")
		cfg.RawV1Alpha1().MachineConfig.MachineKernel = &v1alpha1.KernelConfig{}
		cfg.RawV1Alpha1().MachineConfig.MachineKernel.KernelModules = node.KernelModules
	}

	if node.NodeAnnotations != nil {
		nonTemplateAnnotations, _ := templating.SplitTemplatedMapItems(node.NodeAnnotations)
		slog.Debug(fmt.Sprintf("setting node annotations for %s", nonTemplateAnnotations))
		cfg.RawV1Alpha1().MachineConfig.MachineNodeAnnotations = nonTemplateAnnotations
	}

	if node.NodeLabels != nil {
		nonTemplateLabels, _ := templating.SplitTemplatedMapItems(node.NodeLabels)
		slog.Debug(fmt.Sprintf("setting node labels to %s", nonTemplateLabels))
		cfg.RawV1Alpha1().MachineConfig.MachineNodeLabels = nonTemplateLabels
	}

	if node.NodeTaints != nil {
		slog.Debug(fmt.Sprintf("setting node taints to %s", node.NodeTaints))
		cfg.RawV1Alpha1().MachineConfig.MachineNodeTaints = node.NodeTaints
	}

	if len(node.MachineFiles) > 0 {
		slog.Debug("setting machine files")
		cfg.RawV1Alpha1().MachineConfig.MachineFiles = node.MachineFiles.GetMFs()
	}

	if node.Schematic != nil && len(node.Schematic.Customization.ExtraKernelArgs) > 0 {
		// Talos doesn't support kernel arguments when using SDBoot
		// see: https://github.com/budimanjojo/talhelper/issues/1000
		// Talos 1.12+ defaults to grubUseUKICmdline=true, which is incompatible with install.extraKernelArgs
		// see: https://github.com/budimanjojo/talhelper/issues/1341
		if !node.MachineSpec.UseUKI && !vc.GrubUseUKICmdlineDefault() {
			slog.Debug("appending schematic kernel args to install kernel args")
			//nolint:staticcheck
			cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallExtraKernelArgs = append(cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallExtraKernelArgs, node.Schematic.Customization.ExtraKernelArgs...)
		}
	}

	if len(node.CertSANs) > 0 {
		slog.Debug("appending extra machine certificate SANs")
		cfg.RawV1Alpha1().MachineConfig.MachineCertSANs = append(cfg.RawV1Alpha1().MachineConfig.MachineCertSANs, node.CertSANs...)
	}

	return cfg
}

// Template supported fields within the config
func templateConfig(node *config.Node, cfg taloscfg.Provider) (taloscfg.Provider, error) {
	// Two separate copies of the configuration are required: one that is
	// continuously updated with rendered values, and one that that contains
	// the constant, raw template values. This ensures consistent results in
	// the event that support for new fields are added, or that they are
	// templated in a different order, which is important to prevent breaking
	// changes.
	renderedConfig := cfg.RawV1Alpha1()
	unchangedConfig := renderedConfig.DeepCopy()

	err := templateConfigField(node.NodeLabels, &renderedConfig.MachineConfig.MachineNodeLabels, unchangedConfig, "node labels")
	if err != nil {
		return nil, err
	}

	err = templateConfigField(node.NodeAnnotations, &renderedConfig.MachineConfig.MachineNodeAnnotations, unchangedConfig, "node annotations")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Given a map key pairs, template the values with the provided config, and apply the updated values to the destination map.
// A message is written to the debug log with the field name and templated values.
func templateConfigField[T any](srcKeyPairs map[string]string, dstKeyPairs *map[string]T, cfg *v1alpha1.Config, fieldName string) error {
	_, templateKeyPairs := templating.SplitTemplatedMapItems(srcKeyPairs)
	if len(templateKeyPairs) == 0 {
		return nil
	}

	renderedKeyPairs, err := templating.RenderMap[T](templateKeyPairs, cfg)
	if err != nil {
		return err
	}

	slog.Debug(fmt.Sprintf("adding %v to %s", renderedKeyPairs, fieldName))
	for key := range renderedKeyPairs {
		(*dstKeyPairs)[key] = renderedKeyPairs[key]
	}

	return nil
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
