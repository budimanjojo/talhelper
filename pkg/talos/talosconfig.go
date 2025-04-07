package talos

import (
	"fmt"
	"log/slog"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
)

func GenerateClientConfigBytes(c *config.TalhelperConfig, input *generate.Input, listNodes bool) ([]byte, error) {
	var endpoints []string
	for _, node := range c.Nodes {
		if node.ControlPlane {
			endpoints = append(endpoints, node.GetIPAddresses()...)
		}
	}
	slog.Debug(fmt.Sprintf("endpoints in talosconfig are set to %s", endpoints))
	input.Options.EndpointList = endpoints

	cfg, err := input.Talosconfig()
	if err != nil {
		return nil, err
	}

	// The talos production recommendations recommend explicitly setting --node flags and no default nodes.
	if listNodes {
		slog.Debug("appending all nodes to nodes in talosconfig")
		for _, node := range c.Nodes {
			cfg.Contexts[cfg.Context].Nodes = append(cfg.Contexts[cfg.Context].Nodes, node.GetIPAddresses()...)
		}
	}

	finalCfg, err := cfg.Bytes()
	if err != nil {
		return nil, err
	}

	return finalCfg, nil
}
