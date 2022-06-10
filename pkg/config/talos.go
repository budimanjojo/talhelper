package config

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/talos-systems/crypto/x509"
	"gopkg.in/yaml.v3"

	talosconfig "github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
)

func reEncodeTalosNodeConfig(cfgFile []byte, cfg *v1alpha1.Config) ([]byte, error) {
	err := yaml.Unmarshal(cfgFile, &cfg)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	err = encoder.Encode(cfg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ParseTalosInput(config TalhelperConfig) (*generate.Input, error) {
	kubernetesVersion := config.k8sVersion()

	versionContract, err := talosconfig.ParseContractFromVersion(config.talosVersion())
	if err != nil {
		return nil, err
	}

	secrets, err := generate.NewSecretsBundle(generate.NewClock(), generate.WithVersionContract(versionContract))
	if err != nil {
		return nil, err
	}

	opts := []generate.GenOption{}

	opts = append(opts, generate.WithVersionContract(versionContract))
	opts = append(opts, generate.WithInstallImage(config.installerURL()))

	if config.AllowSchedulingOnMasters {
		opts = append(opts, generate.WithAllowSchedulingOnMasters(config.AllowSchedulingOnMasters))
	}

	if config.CNIConfig.Name != "" {
		opts = append(opts, generate.WithClusterCNIConfig(&v1alpha1.CNIConfig{CNIName: config.CNIConfig.Name, CNIUrls: config.CNIConfig.Urls}))
	}

	if config.Domain != "" {
		opts = append(opts, generate.WithDNSDomain(config.Domain))
	}

	input, err := generate.NewInput(config.ClusterName, config.Endpoint, kubernetesVersion, secrets, opts[:]...)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func createTalosClusterConfig(node nodes, config TalhelperConfig, input *generate.Input) (CfgFile []byte, err error) {
	var cfg *v1alpha1.Config

	var patch []map[string]interface{}
	var iPatch map[string]interface{}

	nodeInput := input
	if node.InstallDisk != "" {
		nodeInput.InstallDisk = node.InstallDisk
	}

	switch node.ControlPlane {
	case true:
		cfg, err = generate.Config(machine.TypeControlPlane, nodeInput)
		if err != nil {
			return nil, err
		}
		patch = config.ControlPlane.ConfigPatches
		iPatch = config.ControlPlane.InlinePatch
	case false:
		cfg, err = generate.Config(machine.TypeWorker, nodeInput)
		if err != nil {
			return nil, err
		}
		patch = config.Worker.ConfigPatches
		iPatch = config.Worker.InlinePatch
	}

	cfg.MachineConfig.MachineNetwork.NetworkHostname = node.Hostname

	if len(node.NetworkInterfaces) != 0 {
		iface := make([]v1alpha1.Device, len(node.NetworkInterfaces))
		for k, v := range node.NetworkInterfaces {
			iface[k].DeviceInterface = v.Interface
			iface[k].DeviceAddresses = v.Addresses
			iface[k].DeviceMTU = v.MTU
			iface[k].DeviceIgnore = v.Ignore
			iface[k].DeviceDHCP = v.DHCP
		}

		for _, v := range iface {
			v := v
			cfg.MachineConfig.MachineNetwork.NetworkInterfaces = append(cfg.MachineConfig.MachineNetwork.NetworkInterfaces, &v)
		}
	}

	marshaledCfg, err := cfg.Bytes()
	if err != nil {
		return nil, err
	}

	marshaledPatch, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	patchedCfg, err := applyPatchFromYaml(marshaledPatch, marshaledCfg)
	if err != nil {
		return nil, err
	}

	if iPatch != nil {
		marshaledIPatch, err := json.Marshal(iPatch)
		if err != nil {
			return nil, err
		}

		finalCfg, err := ApplyInlinePatchFromYaml(marshaledIPatch, patchedCfg)
		if err != nil {
			return nil, err
		}

		return finalCfg, nil
	}
	return patchedCfg, nil
}

func createTalosClientConfig(config TalhelperConfig, input *generate.Input, machineCert *x509.PEMEncodedCertificateAndKey) ([]byte, error) {
	options := generate.DefaultGenOptions()

	var endpointList []string
	for _, node := range config.Nodes {
		if node.ControlPlane {
			endpointList = append(endpointList, node.IPAddress)
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

	clientCfg, err := generate.Talosconfig(input, generate.WithEndpointList(endpointList))
	if err != nil {
		return nil, err
	}

	marshaledClientCfg, err := clientCfg.Bytes()
	if err != nil {
		return nil, err
	}

	return marshaledClientCfg, nil
}
