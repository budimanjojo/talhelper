package talos

import (
	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/decrypt"
	"github.com/budimanjojo/talhelper/pkg/substitute"
	tconfig "github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/generate/secrets"
	"gopkg.in/yaml.v3"
)

// NewClusterInput takes `Talhelper` config and path to encrypted `secretFile` and
// returns Talos `generate.Input`. It also returns an error, if any.
func NewClusterInput(c *config.TalhelperConfig, secretFile string) (*generate.Input, error) {
	kubernetesVersion := c.GetK8sVersion()

	versionContract, err := tconfig.ParseContractFromVersion(c.GetTalosVersion())
	if err != nil {
		return nil, err
	}

	var sb *secrets.Bundle

	if secretFile != "" {
		decrypted, err := decrypt.DecryptYamlWithSops(secretFile)
		if err != nil {
			return nil, err
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
		sb, err = NewSecretBundle(secrets.NewClock(), *versionContract)
		if err != nil {
			return nil, err
		}
	}

	opts := parseOptions(c, versionContract, sb)

	input, err := generate.NewInput(c.ClusterName, c.Endpoint, kubernetesVersion, opts...)
	if err != nil {
		return nil, err
	}

	input.PodNet = c.GetClusterPodNets()
	input.ServiceNet = c.GetClusterSvcNets()
	input.AdditionalMachineCertSANs = c.AdditionalMachineCertSans
	input.AdditionalSubjectAltNames = c.AdditionalApiServerCertSans

	return input, nil
}

// parseOptions takes `TalhelperConfig` and returns slice of Talos `generate.GenOption`
// compatible with the specified `versionContract`.
func parseOptions(c *config.TalhelperConfig, versionContract *tconfig.VersionContract, sb *secrets.Bundle) []generate.Option {
	opts := []generate.Option{}

	opts = append(opts, generate.WithVersionContract(versionContract))
	opts = append(opts, generate.WithSecretsBundle(sb))
	opts = append(opts, generate.WithInstallImage("ghcr.io/siderolabs/installer:"+c.GetTalosVersion()))

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
