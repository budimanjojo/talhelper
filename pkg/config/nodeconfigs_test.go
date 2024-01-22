package config

import (
	"reflect"
	"testing"

	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
)

func TestOverrideNodeConfigs(t *testing.T) {
	globalCfg := NodeConfigs{
		NodeLabels: map[string]string{
			"testkey": "testValue",
		},
		Schematic: &schematic.Schematic{
			Customization: schematic.Customization{
				ExtraKernelArgs: []string{"enable=1"},
			},
		},
	}

	node := Node{
		Hostname:     "test-host",
		IPAddress:    "123.456.789.1",
		InstallDisk:  "/dev/test",
		ControlPlane: true,
		NodeConfigs: NodeConfigs{
			NodeLabels: map[string]string{
				"testkey": "overwritten",
			},
			MachineDisks: []*v1alpha1.MachineDisk{
				{
					DeviceName: "/dev/sda",
					DiskPartitions: []*v1alpha1.DiskPartition{
						{
							DiskSize:       v1alpha1.DiskSize(1),
							DiskMountPoint: "/hello",
						},
					},
				},
			},
		},
	}

	expectedNode := Node{
		Hostname:     "test-host",
		IPAddress:    "123.456.789.1",
		InstallDisk:  "/dev/test",
		ControlPlane: true,
		NodeConfigs: NodeConfigs{
			NodeLabels: map[string]string{
				"testkey": "overwritten",
			},
			MachineDisks: []*v1alpha1.MachineDisk{
				{
					DeviceName: "/dev/sda",
					DiskPartitions: []*v1alpha1.DiskPartition{
						{
							DiskSize:       v1alpha1.DiskSize(1),
							DiskMountPoint: "/hello",
						},
					},
				},
			},
			Schematic: &schematic.Schematic{
				Customization: schematic.Customization{
					ExtraKernelArgs: []string{"enable=1"},
				},
			},
		},
	}

	node.OverrideGlobalCfg(globalCfg)

	if !reflect.DeepEqual(node, expectedNode) {
		t.Errorf("got:\n%v\nwant:\n%v", node, expectedNode)
	}
}
