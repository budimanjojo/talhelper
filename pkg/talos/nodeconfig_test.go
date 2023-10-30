package talos

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	"gopkg.in/yaml.v3"
)

func TestGenerateNodeConfig(t *testing.T) {
	data := []byte(`clusterName: test
talosVersion: v1.5.4
endpoint: https://1.1.1.1:6443
nodes:
  - hostname: node1
    controlPlane: true
    installDisk: /dev/sda
    disableSearchDomain: true
    extensions:
      - image: ghcr.io/siderolabs/tailscale:1.44.0
    machineFiles:
      - content: TS_AUTHKEY=123456
        permissions: 0o644
        path: /var/etc/tailscale/auth.env
        op: create
    networkInterfaces:
      - interface: eth0
        bond:
          deviceSelectors:
            - hardwareAddr: "00:50:56:*"
    kernelModules:
      - name: br_netfilter
        parameters:
          - nf_conntrack_max=131072
    schematic:
      customization:
        extraKernelArgs:
          - hello
          - hihi
        systemExtensions:
          officialExtensions:
            - siderolabs/amd-ucode
            - siderolabs/nvidia-fabricmanager
  - hostname: node2
    controlPlane: false
    installDiskSelector:
      size: 4KB
      model: WDC*
      name: /sys/block/sda/device/name
      busPath: /pci0000:00/0000:00:17.0/ata1/host0/target0:0:0/0:0:0:0
    machineDisks:
      - device: /dev/disk/by-id/ata-CT500MX500SSD1_2149E5EC1D9D
        partitions:
          - mountpoint: /var/mnt/sata
    nodeLabels:
      rack: rack1a
      zone: us-east-1a
    talosImageURL: factory.talos.dev/installer/e9c7ef96884d4fbc8c0a1304ccca4bb0287d766a8b4125997cb9dbe84262144e
    schematic:
      customization:
        extraKernelArgs:
          - hello
          - hihi
        systemExtensions:
          officialExtensions:
            - siderolabs/amd-ucode
            - siderolabs/nvidia-fabricmanager
    nameservers: [1.1.1.1, 8.8.8.8]`)

	var m config.TalhelperConfig
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	input, err := NewClusterInput(&m, "")
	if err != nil {
		t.Fatal(err)
	}

	cp, err := GenerateNodeConfig(&m.Nodes[0], input)
	if err != nil {
		t.Fatal(err)
	}

	w, err := GenerateNodeConfig(&m.Nodes[1], input)
	if err != nil {
		t.Fatal(err)
	}

	expectedNode1Type := "controlplane"
	expectedNode1Hostname := "node1"
	expectedNode1InstallDisk := "/dev/sda"
	expectedNode1DisableSearchDomain := true
	expectedNode1Extensions := []v1alpha1.InstallExtensionConfig{
		{
			ExtensionImage: "ghcr.io/siderolabs/tailscale:1.44.0",
		},
	}
	expectedNode1MachineFiles := []*v1alpha1.MachineFile{
		{
			FileContent:     "TS_AUTHKEY=123456",
			FilePermissions: v1alpha1.FileMode(0o644),
			FilePath:        "/var/etc/tailscale/auth.env",
			FileOp:          "create",
		},
	}
	expectedNode1NetworkInterfaces := v1alpha1.NetworkDeviceList{
		{
			DeviceInterface: "eth0",
			DeviceBond: &v1alpha1.Bond{
				BondDeviceSelectors: []v1alpha1.NetworkDeviceSelector{
					{
						NetworkDeviceHardwareAddress: "00:50:56:*",
					},
				},
			},
		},
	}
	expectedNode1KernelModules := &v1alpha1.KernelConfig{
		KernelModules: []*v1alpha1.KernelModuleConfig{
			{
				ModuleName:       "br_netfilter",
				ModuleParameters: []string{"nf_conntrack_max=131072"},
			},
		},
	}
	expectedNode1InstallImage := "factory.talos.dev/installer/647a0a54bff662aa12051bc0312097f29d3562107d8e6a8e87ab85b643e25bc0:v1.5.4"
	expectedNode2Type := "worker"
	expectedNode2InstallDiskSelector := &v1alpha1.InstallDiskSelector{
		Size: &v1alpha1.InstallDiskSizeMatcher{
			MatchData: v1alpha1.InstallDiskSizeMatchData{
				Size: 4000,
				Op:   "",
			},
		},
		Model:   "WDC*",
		Name:    "/sys/block/sda/device/name",
		BusPath: "/pci0000:00/0000:00:17.0/ata1/host0/target0:0:0/0:0:0:0",
	}
	expectedNode2MachineDisks := []*v1alpha1.MachineDisk{
		{
			DeviceName: "/dev/disk/by-id/ata-CT500MX500SSD1_2149E5EC1D9D",
			DiskPartitions: []*v1alpha1.DiskPartition{
				{
					DiskMountPoint: "/var/mnt/sata",
				},
			},
		},
	}
	expectedNode2NodeLabels := map[string]string{"rack": "rack1a", "zone": "us-east-1a"}
	expectedNode2InstallImage := "factory.talos.dev/installer/e9c7ef96884d4fbc8c0a1304ccca4bb0287d766a8b4125997cb9dbe84262144e:v1.5.4"
	expectedNode2Nameservers := []string{"1.1.1.1", "8.8.8.8"}

	cpCfg := cp.RawV1Alpha1().MachineConfig
	wCfg := w.RawV1Alpha1().MachineConfig

	compare(cpCfg.MachineType, expectedNode1Type, t)
	compare(cpCfg.MachineNetwork.Hostname(), expectedNode1Hostname, t)
	compare(cpCfg.MachineInstall.InstallDisk, expectedNode1InstallDisk, t)
	compare(cpCfg.MachineNetwork.DisableSearchDomain(), expectedNode1DisableSearchDomain, t)
	compare(cpCfg.MachineInstall.InstallExtensions, expectedNode1Extensions, t)
	compare(cpCfg.MachineFiles, expectedNode1MachineFiles, t)
	compare(cpCfg.MachineNetwork.NetworkInterfaces, expectedNode1NetworkInterfaces, t)
	compare(cpCfg.MachineKernel, expectedNode1KernelModules, t)
	compare(cpCfg.MachineInstall.InstallImage, expectedNode1InstallImage, t)
	compare(wCfg.MachineType, expectedNode2Type, t)
	compare(wCfg.MachineInstall.InstallDiskSelector, expectedNode2InstallDiskSelector, t)
	compare(wCfg.MachineDisks, expectedNode2MachineDisks, t)
	compare(wCfg.MachineNodeLabels, expectedNode2NodeLabels, t)
	compare(wCfg.MachineInstall.InstallImage, expectedNode2InstallImage, t)
	compare(wCfg.MachineNetwork.NameServers, expectedNode2Nameservers, t)
}

func compare(got, want any, t *testing.T) {
	// Indicate this function is a helper and we're not interested in line numbers coming from it
	t.Helper()
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("\ngot type of %s, want type of %s", reflect.TypeOf(got), reflect.TypeOf(want))
		return
	}
	switch got.(type) {
	case string, int, bool, float64, float32:
		if got != want {
			t.Errorf("\ngot %s\nwant %s", got, want)
		}
	case map[string]string:
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot %s\nwant %s", got, want)
		}
	default:
		g, err := json.Marshal(got)
		if err != nil {
			t.Errorf("error encoding %v to json", got)
		}
		w, err := json.Marshal(want)
		if err != nil {
			t.Errorf("error encoding %v to json", want)
		}
		if !reflect.DeepEqual(g, w) {
			t.Errorf("\ngot %s\nwant %s", g, w)
		}
	}
}
