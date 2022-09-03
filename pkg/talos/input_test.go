package talos

import (
	"testing"

	"github.com/budimanjojo/talhelper/pkg/config"
	"gopkg.in/yaml.v3"
)

func TestNewClusterInput(t *testing.T) {
	data := []byte(`clusterName: test
talosVersion: 1.0
kubernetesVersion: v1.24.1
endpoint: https://1.1.1.1:6443`)

	var m config.TalhelperConfig

	err := yaml.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	result, err := NewClusterInput(&m, "")
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
