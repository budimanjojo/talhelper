package talos

import (
	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/types/runtime/extensions"
)

func GenerateExtensionServicesConfigBytes(esCfgs []*config.ExtensionService) ([]byte, error) {
	var result [][]byte

	exts, err := GenerateNodeExtensionServiceConfig(esCfgs)
	if err != nil {
		return nil, err
	}

	for _, ext := range exts {
		extByte, err := marshalYaml(ext)
		if err != nil {
			return nil, err
		}

		result = append(result, extByte)
	}

	return CombineYamlBytes(result), nil
}

func GenerateNodeExtensionServiceConfig(esCfgs []*config.ExtensionService) ([]*extensions.ServiceConfigV1Alpha1, error) {
	var result []*extensions.ServiceConfigV1Alpha1

	for _, v := range esCfgs {
		esc := extensions.NewServicesConfigV1Alpha1()
		esc.ServiceName = v.Name
		esc.ServiceConfigFiles = v.ConfigFiles
		esc.ServiceEnvironment = v.Environment

		if _, err := esc.Validate(nil); err != nil {
			return nil, err
		}

		result = append(result, esc)
	}

	return result, nil
}
