package talos

import (
	"time"

	"github.com/talos-systems/crypto/x509"

	"github.com/budimanjojo/talhelper/pkg/config"
	tconfig "github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
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

func GenerateNodeConfigBytes(node *config.Nodes, input *generate.Input) ([]byte, error) {
	cfg, err := generateNodeConfig(node, input)
	if err != nil {
		return nil, err
	}
	return cfg.EncodeBytes()
}

func generateNodeConfig(node *config.Nodes, input *generate.Input) (*v1alpha1.Config, error) {
	var c *v1alpha1.Config
	var err error

	nodeInput, err := patchNodeInput(node, input)
	if err != nil {
		return nil, err
	}

	switch node.ControlPlane {
	case true:
		c, err = generate.Config(machine.TypeControlPlane, nodeInput)
		if err != nil {
			return nil, err
		}
	case false:
		c, err = generate.Config(machine.TypeWorker, nodeInput)
		if err != nil {
			return nil, err
		}
	}

	cfg := applyNodeOverride(node, c)

	return cfg, nil
}

func GenerateClientConfigBytes(c *config.TalhelperConfig, input *generate.Input, machineCert *x509.PEMEncodedCertificateAndKey) ([]byte, error) {
	options := generate.DefaultGenOptions()

	var endpoints []string
	for _, node := range c.Nodes {
		if node.ControlPlane {
			endpoints = append(endpoints, node.IPAddress)
		}
	}

	// make sure ca in talosconfig match machine.ca.crt in machine config
	if string(input.Certs.OS.Crt) != string(machineCert.Crt) {
		input.Certs.OS = machineCert

		adminCert, err := generate.NewAdminCertificateAndKey(time.Now(), machineCert, options.Roles, 87600*time.Hour)
		if err != nil {
			return nil, err
		}

		input.Certs.Admin = adminCert
	}

	cfg, err := generate.Talosconfig(input, generate.WithEndpointList(endpoints))
	if err != nil {
		return nil, err
	}

	finalCfg, err := cfg.Bytes()
	if err != nil {
		return nil, err
	}

	return finalCfg, nil
}

func applyNodeOverride(node *config.Nodes, cfg *v1alpha1.Config) (*v1alpha1.Config) {
	cfg.MachineConfig.MachineNetwork.NetworkHostname = node.Hostname

	if len(node.Nameservers) != 0 {
		cfg.MachineConfig.MachineNetwork.NameServers = node.Nameservers
	}

	if len(node.NetworkInterfaces) != 0 {
		iface := make([]v1alpha1.Device, len(node.NetworkInterfaces))
		for k, v := range node.NetworkInterfaces {
			iface[k].DeviceInterface = v.Interface
			iface[k].DeviceAddresses = v.Addresses
			iface[k].DeviceMTU = v.MTU
			iface[k].DeviceIgnore = v.Ignore
			iface[k].DeviceDHCP = v.DHCP
			var route []v1alpha1.Route
			if len(v.Routes) != 0 {
				route = make([]v1alpha1.Route, len(v.Routes))
				for k2, v2 := range v.Routes {
					route[k2].RouteGateway = v2.Gateway
					route[k2].RouteNetwork = v2.Network
					route[k2].RouteMetric = v2.Metric
					route[k2].RouteSource = v2.Source
				}
			}
			for _, v := range route {
				v := v
				iface[k].DeviceRoutes = append(iface[k].DeviceRoutes, &v)
			}
		}

		for _, v := range iface {
			v := v
			cfg.MachineConfig.MachineNetwork.NetworkInterfaces = append(cfg.MachineConfig.MachineNetwork.NetworkInterfaces, &v)
		}
	}

	return cfg
}

func patchNodeInput(node *config.Nodes, input *generate.Input) (*generate.Input, error) {
	nodeInput := input
	if node.InstallDisk != "" {
		nodeInput.InstallDisk = node.InstallDisk
	}

	return nodeInput, nil
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
