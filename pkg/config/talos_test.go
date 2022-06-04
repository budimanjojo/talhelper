package config

import (
	"testing"

	"github.com/talos-systems/talos/pkg/machinery/config/configloader"
	"sigs.k8s.io/yaml"
)

func TestParseTalosInput(t *testing.T) {
	data := []byte(`clusterName: test
talosVersion: v1.0
kubernetesVersion: 1.24.1
endpoint: https://1.1.1.1:6443`)

	var m TalhelperConfig

	err := yaml.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ParseTalosInput(m)
	if err != nil {
		t.Fatal(err)
	}

	expectedClusterName := "test"
	expectedTalosMajVersion := 1
	expectedTalosMinVersion := 0
	expectedK8sVersion := "1.24.1"
	expectedEndpoint := "https://1.1.1.1:6443"

	if result.ClusterName != expectedClusterName {
		t.Errorf("got %s, want %s", result.ClusterName, expectedClusterName)
	}

	if result.VersionContract.Major != expectedTalosMajVersion {
		t.Errorf("got %d, want %d", result.VersionContract.Major, expectedTalosMajVersion)
	}

	if result.VersionContract.Minor != expectedTalosMinVersion {
		t.Errorf("got %d, want %d", result.VersionContract.Minor, expectedTalosMinVersion)
	}

	if result.KubernetesVersion != expectedK8sVersion {
		t.Errorf("got %s, want %s", result.KubernetesVersion, expectedK8sVersion)
	}

	if result.ControlPlaneEndpoint != expectedEndpoint {
		t.Errorf("got %s, want %s", result.ControlPlaneEndpoint, expectedEndpoint)
	}
}

func TestCreateTalosClusterConfig(t *testing.T) {
	data := []byte(`clusterName: test
endpoint: https://1.1.1.1:6443
nodes:
  - hostname: node1
    controlPlane: true
    installDisk: /dev/sda
  - hostname: node2
    controlPlane: false`)

	var m TalhelperConfig

	err := yaml.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	input, err := ParseTalosInput(m)
	if err != nil {
		t.Fatal(err)
	}

	cp, err := createTalosClusterConfig(m.Nodes[0], m, input)
	if err != nil {
		t.Fatal(err)
	}

	w, err := createTalosClusterConfig(m.Nodes[1], m, input)
	if err != nil {
		t.Fatal(err)
	}

	cpCfg, err := configloader.NewFromBytes(cp)
	if err != nil {
		t.Fatal(err)
	}

	wCfg, err := configloader.NewFromBytes(w)
	if err != nil {
		t.Fatal(err)
	}

	expectedNode1Type := "controlplane"
	expectedNode1Hostname := "node1"
	expectedNode1Disk := "/dev/sda"
	expectedNode2Type := "worker"

	if cpCfg.Machine().Type().String() != expectedNode1Type {
		t.Errorf("got %s, want %s", cpCfg.Machine().Type().String(), expectedNode1Type)
	}

	if cpCfg.Machine().Network().Hostname() != expectedNode1Hostname {
		t.Errorf("got %s, want %s", cpCfg.Machine().Network().Hostname(), expectedNode1Hostname)
	}

	if node1Disk, _ := cpCfg.Machine().Install().Disk(); node1Disk != expectedNode1Disk {
		t.Errorf("got %s, want %s", node1Disk, expectedNode1Disk)
	}

	if wCfg.Machine().Type().String() != expectedNode2Type {
		t.Errorf("got %s, want %s", cpCfg.Machine().Type().String(), expectedNode1Type)
	}
}
