package talos

import (
	"testing"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/cel"
	"github.com/siderolabs/talos/pkg/machinery/cel/celenv"
	"github.com/siderolabs/talos/pkg/machinery/config/types/block"
	"gopkg.in/yaml.v3"
)

func TestGenerateNodeVolumeConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    volumes:
      - name: EPHEMERAL
        provisioning:
          diskSelector:
            match: disk.transport == "nvme"
          maxSize: 50GiB
      - name: IMAGECACHE
        provisioning:
          diskSelector:
            match: disk.size > 120u * GB && disk.size < 1u * TB`)
	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	expectedVolume1Name := "EPHEMERAL"
	expectedVolume1Provisioning := block.ProvisioningSpec{
		DiskSelectorSpec: block.DiskSelector{
			Match: cel.MustExpression(cel.ParseBooleanExpression(`disk.transport == "nvme"`, celenv.DiskLocator())),
		},
		ProvisioningMaxSize: block.MustByteSize("50GiB"),
	}
	expectedVolume2Name := "IMAGECACHE"
	expectedVolume2Provisioning := block.ProvisioningSpec{
		DiskSelectorSpec: block.DiskSelector{
			Match: cel.MustExpression(cel.ParseBooleanExpression(`disk.size > 120u * GB && disk.size < 1u * TB`, celenv.DiskLocator())),
		},
	}

	result, err := GenerateVolumeConfig(m.Nodes[0].Volumes, "metal")
	if err != nil {
		t.Fatal(err)
	}

	compare(result[0].Name(), expectedVolume1Name, t)
	compare(result[0].ProvisioningSpec, expectedVolume1Provisioning, t)
	compare(result[1].Name(), expectedVolume2Name, t)
	compare(result[1].ProvisioningSpec, expectedVolume2Provisioning, t)
}
