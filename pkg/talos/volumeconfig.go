package talos

import (
	"fmt"
	"slices"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/types/block"
)

func GenerateVolumeConfigBytes(cfgs []*config.Volume, mode string) ([]byte, error) {
	var result [][]byte

	vcs, err := GenerateVolumeConfig(cfgs, mode)
	if err != nil {
		return nil, err
	}

	for _, vc := range vcs {
		vcByte, err := marshalYaml(vc)
		if err != nil {
			return nil, err
		}

		result = append(result, vcByte)
	}

	return CombineYamlBytes(result), nil
}

func GenerateVolumeConfig(cfgs []*config.Volume, mode string) ([]*block.VolumeConfigV1Alpha1, error) {
	var (
		// I suppose we shouldn't allow same volume names?
		names  []string
		result []*block.VolumeConfigV1Alpha1
	)

	m, err := parseMode(mode)
	if err != nil {
		return nil, err
	}

	for _, v := range cfgs {
		if slices.Index(names, v.Name) != -1 {
			return nil, fmt.Errorf("duplicated volume config name found: %s", v.Name)
		}
		names = append(names, v.Name)
		vc := block.NewVolumeConfigV1Alpha1()
		vc.MetaName = v.Name
		vc.ProvisioningSpec = v.Provisioning
		vc.EncryptionSpec = v.Encryption

		if _, err := vc.Validate(m); err != nil {
			return nil, err
		}

		result = append(result, vc)
	}

	return result, nil
}
