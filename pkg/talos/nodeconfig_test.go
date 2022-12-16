package talos

import (
	"testing"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	"gopkg.in/yaml.v3"
)

func TestGenerateNodeConfigBytes(t *testing.T) {
	data := []byte(`clusterName: test
endpoint: https://1.1.1.1:6443
nodes:
  - hostname: node1
    controlPlane: true
    installDisk: /dev/sda
  - hostname: node2
    controlPlane: false`)

	var m config.TalhelperConfig
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	input, err := NewClusterInput(&m, "")
	if err != nil {
		t.Fatal(err)
	}

	cp, err := GenerateNodeConfigBytes(&m.Nodes[0], input)
	if err != nil {
		t.Fatal(err)
	}

	w, err := GenerateNodeConfigBytes(&m.Nodes[1], input)
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
