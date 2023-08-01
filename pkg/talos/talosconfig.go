package talos

import (
	"time"

	"github.com/siderolabs/crypto/x509"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/generate/secrets"
)

func GenerateClientConfigBytes(c *config.TalhelperConfig, input *generate.Input, machineCert *x509.PEMEncodedCertificateAndKey) ([]byte, error) {
	options := generate.DefaultOptions()

	var endpoints []string
	for _, node := range c.Nodes {
		if node.ControlPlane {
			endpoints = append(endpoints, node.IPAddress)
		}
	}
	input.Options.EndpointList = endpoints

	// make sure ca in talosconfig match machine.ca.crt in machine config
	if string(input.Options.SecretsBundle.Certs.OS.Crt) != string(machineCert.Crt) {
		input.Options.SecretsBundle.Certs.OS = machineCert

		adminCert, err := secrets.NewAdminCertificateAndKey(time.Now(), machineCert, options.Roles, 87600*time.Hour)
		if err != nil {
			return nil, err
		}

		input.Options.SecretsBundle.Certs.Admin = adminCert
	}

	cfg, err := input.Talosconfig()
	if err != nil {
		return nil, err
	}

	finalCfg, err := cfg.Bytes()
	if err != nil {
		return nil, err
	}

	return finalCfg, nil
}
