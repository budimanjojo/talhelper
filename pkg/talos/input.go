package talos

import (
	"github.com/budimanjojo/talhelper/pkg/config"
	tconfig "github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
)

func NewClusterInput(c *config.TalhelperConfig) (*generate.Input, error) {
	kubernetesVersion := c.GetK8sVersion()

	versionContract, err := tconfig.ParseContractFromVersion(c.GetTalosVersion())
	if err != nil {
		return nil, err
	}

	secrets, err := generate.NewSecretsBundle(generate.NewClock(), generate.WithVersionContract(versionContract))
	if err != nil {
		return nil, err
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

	if c.AllowSchedulingOnMasters {
		opts = append(opts, generate.WithAllowSchedulingOnMasters(c.AllowSchedulingOnMasters))
	}

	if c.CNIConfig.Name != "" {
		opts = append(opts, generate.WithClusterCNIConfig(&v1alpha1.CNIConfig{CNIName: c.CNIConfig.Name, CNIUrls: c.CNIConfig.Urls}))
	}

	if c.Domain != "" {
		opts = append(opts, generate.WithDNSDomain(c.Domain))
	}

	return opts
}
