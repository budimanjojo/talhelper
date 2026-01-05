package talos

import (
	"fmt"
	"slices"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/types/block"
)

func GenerateUserVolumeConfigBytes(cfgs []*config.UserVolume, mode string) ([]byte, error) {
	var result [][]byte

	uvcs, err := GenerateUserVolumeConfig(cfgs, mode)
	if err != nil {
		return nil, err
	}

	for _, uvc := range uvcs {
		uvcByte, err := marshalYaml(uvc)
		if err != nil {
			return nil, err
		}

		result = append(result, uvcByte)
	}

	return CombineYamlBytes(result), nil
}

func GenerateUserVolumeConfig(cfgs []*config.UserVolume, mode string) ([]*block.UserVolumeConfigV1Alpha1, error) {
	var (
		// I suppose we shouldn't allow same volume names?
		names  []string
		result []*block.UserVolumeConfigV1Alpha1
	)

	m, err := parseMode(mode)
	if err != nil {
		return nil, err
	}

	for _, uv := range cfgs {
		if slices.Index(names, uv.Name) != -1 {
			return nil, fmt.Errorf("duplicated user volume config name found: %s", uv.Name)
		}
		names = append(names, uv.Name)
		uvc := block.NewUserVolumeConfigV1Alpha1()
		uvc.MetaName = uv.Name
		uvc.VolumeType = uv.VolumeType
		uvc.ProvisioningSpec = uv.Provisioning
		uvc.FilesystemSpec = uv.Filesystem
		uvc.EncryptionSpec = uv.Encryption

		if _, err := uvc.Validate(m); err != nil {
			return nil, err
		}

		result = append(result, uvc)
	}

	return result, nil
}
