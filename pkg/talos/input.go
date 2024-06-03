package talos

import (
	"fmt"
	"log/slog"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/decrypt"
	"github.com/budimanjojo/talhelper/v3/pkg/substitute"
	tconfig "github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/generate/secrets"
	"gopkg.in/yaml.v3"
)

// NewClusterInput takes `Talhelper` config and path to encrypted `secretFile` and
// returns Talos `generate.Input`. It also returns an error, if any.
func NewClusterInput(c *config.TalhelperConfig, secretFile string, mode string) (*generate.Input, error) {
	kubernetesVersion := c.GetK8sVersion()

	versionContract, err := tconfig.ParseContractFromVersion(c.GetTalosVersion())
	if err != nil {
		return nil, err
	}

	var sb *secrets.Bundle

	if secretFile != "" {
		slog.Debug(fmt.Sprintf("using secret file %s", secretFile))
		decrypted, err := decrypt.DecryptYamlWithSops(secretFile)
		if err != nil {
			return nil, err
		}
		if len(decrypted) == 0 {
			return nil, fmt.Errorf("secret file %s is empty", secretFile)
		}

		decrypted, err = substitute.SubstituteEnvFromByte(decrypted)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(decrypted, &sb)
		if err != nil {
			return nil, err
		}
		sb.Clock = secrets.NewClock()
	} else {
		slog.Debug("generating new secret file because secret file is not found")
		sb, err = NewSecretBundle(secrets.NewClock(), *versionContract)
		if err != nil {
			return nil, err
		}
	}

	opts := parseOptions(c, versionContract, sb, mode)

	slog.Debug("generating input file", "clusterName", c.ClusterName, "endpoint", c.Endpoint, "kubernetesVersion", kubernetesVersion)
	input, err := generate.NewInput(c.ClusterName, c.Endpoint, kubernetesVersion, opts...)
	if err != nil {
		return nil, err
	}

	slog.Debug(fmt.Sprintf("setting input pod network to %s", c.GetClusterPodNets()))
	input.PodNet = c.GetClusterPodNets()
	slog.Debug(fmt.Sprintf("setting input service network to %s", c.GetClusterSvcNets()))
	input.ServiceNet = c.GetClusterSvcNets()
	slog.Debug(fmt.Sprintf("setting input additional machine cert SANs to %s", c.AdditionalMachineCertSans))
	input.AdditionalMachineCertSANs = c.AdditionalMachineCertSans
	slog.Debug(fmt.Sprintf("setting input additional subject alt names to %s", c.AdditionalApiServerCertSans))
	input.AdditionalSubjectAltNames = c.AdditionalApiServerCertSans

	return input, nil
}

// parseOptions takes `TalhelperConfig` and returns slice of Talos `generate.GenOption`
// compatible with the specified `versionContract`.
func parseOptions(c *config.TalhelperConfig, versionContract *tconfig.VersionContract, sb *secrets.Bundle, mode string) []generate.Option {
	opts := []generate.Option{}

	opts = append(opts, generate.WithVersionContract(versionContract))
	opts = append(opts, generate.WithSecretsBundle(sb))
	opts = append(opts, generate.WithInstallImage("ghcr.io/siderolabs/installer:"+c.GetTalosVersion()))

	m, _ := parseMode(mode)
	if m == modeContainer {
		opts = append(opts, generate.WithHostDNSForwardKubeDNSToHost(true))
	}

	if c.AllowSchedulingOnMasters || c.AllowSchedulingOnControlPlanes {
		opts = append(opts, generate.WithAllowSchedulingOnControlPlanes(true))
	}

	if c.CNIConfig != nil {
		opts = append(opts, generate.WithClusterCNIConfig(c.CNIConfig))
	}

	if c.Domain != "" {
		opts = append(opts, generate.WithDNSDomain(c.Domain))
	}

	return opts
}
