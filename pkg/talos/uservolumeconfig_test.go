package talos

import (
	"testing"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/cel"
	"github.com/siderolabs/talos/pkg/machinery/cel/celenv"
	"github.com/siderolabs/talos/pkg/machinery/config/types/block"
	blocktype "github.com/siderolabs/talos/pkg/machinery/resources/block"
	"gopkg.in/yaml.v3"
)

func TestGenerateNodeUserVolumeConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    userVolumes:
      - name: ceph-data
        provisioning:
          diskSelector:
            match: disk.transport == "nvme"
          maxSize: 50GiB
        filesystem:
          type: xfs
      - name: ceph-data2
        provisioning:
          diskSelector:
            match: disk.size > 120u * GB && disk.size < 1u * TB
          minSize: 1GiB`)
	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	expectedVolume1Name := "ceph-data"
	expectedVolume1Provisioning := block.ProvisioningSpec{
		DiskSelectorSpec: block.DiskSelector{
			Match: cel.MustExpression(cel.ParseBooleanExpression(`disk.transport == "nvme"`, celenv.DiskLocator())),
		},
		ProvisioningMaxSize: block.MustByteSize("50GiB"),
	}
	expectedVolume1Filesystem := block.FilesystemSpec{
		FilesystemType: blocktype.FilesystemTypeXFS,
	}
	expectedVolume2Name := "ceph-data2"
	expectedVolume2Provisioning := block.ProvisioningSpec{
		DiskSelectorSpec: block.DiskSelector{
			Match: cel.MustExpression(cel.ParseBooleanExpression(`disk.size > 120u * GB && disk.size < 1u * TB`, celenv.DiskLocator())),
		},
		ProvisioningMinSize: block.MustByteSize("1GiB"),
	}

	result, err := GenerateUserVolumeConfig(m.Nodes[0].UserVolumes, "metal")
	if err != nil {
		t.Fatal(err)
	}

	compare(result[0].Name(), expectedVolume1Name, t)
	compare(result[0].ProvisioningSpec, expectedVolume1Provisioning, t)
	compare(result[0].FilesystemSpec, expectedVolume1Filesystem, t)
	compare(result[1].Name(), expectedVolume2Name, t)
	compare(result[1].ProvisioningSpec, expectedVolume2Provisioning, t)
}
