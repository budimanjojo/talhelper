package talos

import (
	"os"

	"github.com/budimanjojo/talhelper/pkg/config"
	tconfig "github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
)

func NewClusterInput(c *config.TalhelperConfig, secretFile string) (*generate.Input, error) {
	kubernetesVersion := c.GetK8sVersion()

	versionContract, err := tconfig.ParseContractFromVersion(c.GetTalosVersion())
	if err != nil {
		return nil, err
	}

	var secrets *generate.SecretsBundle

	if secretFile != "" {
		secrets, err = NewSecretBundle(generate.NewClock(), generate.WithVersionContract(versionContract), generate.WithSecrets(secretFile))
		if err != nil {
			return nil, err
		}

		err = os.Remove(secretFile)
		if err != nil {
			return nil, err
		}
	} else {
		secrets, err = NewSecretBundle(generate.NewClock(), generate.WithVersionContract(versionContract))
		if err != nil {
			return nil, err
		}
	}

	opts := parseOptions(c, versionContract)

	input, err := generate.NewInput(c.ClusterName, c.Endpoint, kubernetesVersion, secrets, opts...)
	if err != nil {
		return nil, err
	}

	input.PodNet = c.GetClusterPodNets()
	input.ServiceNet = c.GetClusterSvcNets()

	return input, nil
}

func parseOptions(c *config.TalhelperConfig, vc *tconfig.VersionContract) []generate.GenOption {
	opts := []generate.GenOption{}

	opts = append(opts, generate.WithVersionContract(vc))
	opts = append(opts, generate.WithInstallImage(c.GetInstallerURL()))

	if c.AllowSchedulingOnMasters || c.AllowSchedulingOnControlPlanes {
		opts = append(opts, generate.WithAllowSchedulingOnControlPlanes(true))
	}

	if c.CNIConfig.Name != "" {
		opts = append(opts, generate.WithClusterCNIConfig(&v1alpha1.CNIConfig{CNIName: c.CNIConfig.Name, CNIUrls: c.CNIConfig.Urls}))
	}

	if c.Domain != "" {
		opts = append(opts, generate.WithDNSDomain(c.Domain))
	}

	return opts
}
