package talos

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	clientconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
)

func GenerateClientConfigBytes(c *config.TalhelperConfig, input *generate.Input, disableNodesSection bool, crtTTL time.Duration) ([]byte, error) {
	var endpoints []string
	for _, node := range c.Nodes {
		if node.ControlPlane {
			endpoints = append(endpoints, node.GetIPAddresses()...)
		}
	}
	slog.Debug(fmt.Sprintf("endpoints in talosconfig are set to %s", endpoints))

	slog.Debug(fmt.Sprintf("generating admin certificate with TTL of %s", crtTTL))
	cert, err := input.Options.SecretsBundle.GenerateTalosAPIClientCertificateWithTTL(input.Options.Roles, crtTTL)
	if err != nil {
		return nil, err
	}

	cfg := clientconfig.NewConfig(input.ClusterName, endpoints, input.Options.SecretsBundle.Certs.OS.Crt, cert)

	// The talos production recommendations recommend explicitly setting --node flags and no default nodes.
	if !disableNodesSection {
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
