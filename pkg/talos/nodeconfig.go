package talos

import (
	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
)

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

