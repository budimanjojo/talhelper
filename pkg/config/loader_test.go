package config

import (
	"reflect"
	"testing"

	"github.com/siderolabs/image-factory/pkg/schematic"
)

func TestLoadAndValidateFromFile(t *testing.T) {
	cfg, err := LoadAndValidateFromFile("testdata/talconfig.yaml", []string{"testdata/env1.yaml", "testdata/env2.yml"})
	if err != nil {
		t.Fatal(err)
	}

	expectedClusterName := "test-cluster"
	expectedNode0 := Node{
		Hostname:     "hostname1",
		IPAddress:    "192.168.200.10",
		ControlPlane: true,
		InstallDisk:  "/dev/sda",
		Schematic: &schematic.Schematic{
			Customization: schematic.Customization{
				SystemExtensions: schematic.SystemExtensions{
					OfficialExtensions: []string{"siderolabs/tailscale"},
				},
			},
		},
	}
	expectedNode1 := Node{
		Hostname:     "hostname2",
		IPAddress:    "192.168.200.11",
		ControlPlane: false,
		InstallDisk:  "/dev/sda",
		Schematic: &schematic.Schematic{
			Customization: schematic.Customization{
				ExtraKernelArgs: []string{"net.ifnames=0"},
			},
		},
	}

	if cfg.ClusterName != expectedClusterName {
		t.Errorf("got %s, want %s", cfg.ClusterName, expectedClusterName)
	}

	if !reflect.DeepEqual(cfg.Nodes[0], expectedNode0) {
		t.Errorf("got:\n%v\nwant:\n%v", cfg.Nodes[0], expectedNode0)
	}

	if !reflect.DeepEqual(cfg.Nodes[1], expectedNode1) {
		t.Errorf("got:\n%v\nwant:\n%v", cfg.Nodes[1], expectedNode1)
	}
}
