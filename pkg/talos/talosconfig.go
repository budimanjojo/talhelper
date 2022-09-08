package talos

import (
	"time"

	"github.com/talos-systems/crypto/x509"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
)

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
