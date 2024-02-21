package talos

import (
	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
)

func GenerateClientConfigBytes(c *config.TalhelperConfig, input *generate.Input) ([]byte, error) {
	var endpoints []string
	for _, node := range c.Nodes {
		if node.ControlPlane {
			endpoints = append(endpoints, node.GetIPAddresses()...)
		}
	}
	input.Options.EndpointList = endpoints

	cfg, err := input.Talosconfig()
	if err != nil {
		return nil, err
	}

	for _, node := range c.Nodes {
		cfg.Contexts[cfg.Context].Nodes = append(cfg.Contexts[cfg.Context].Nodes, node.GetIPAddresses()...)
	}

	finalCfg, err := cfg.Bytes()
	if err != nil {
		return nil, err
	}

	return finalCfg, nil
}
