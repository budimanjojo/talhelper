package patcher

import (
	"os"
	"strings"

	"github.com/budimanjojo/talhelper/v3/pkg/decrypt"
	"github.com/budimanjojo/talhelper/v3/pkg/substitute"
	"github.com/siderolabs/talos/pkg/machinery/config/configpatcher"
	"gopkg.in/yaml.v3"
)

// YAMLInlinePatcher applies JSON7396 patches into target and returns it.
// It also returns an error, if any.
func YAMLInlinePatcher(patch interface{}, target []byte) ([]byte, error) {
	p, err := yaml.Marshal(patch)
	if err != nil {
		return nil, err
	}

	out, err := JSON7396FromYAML(p, target)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// YAMLPatcher applies JSON6902 patches into target and returns it.
// It also returns an error, if any.
func YAMLPatcher(patch interface{}, target []byte) ([]byte, error) {
	p, err := yaml.Marshal(patch)
	if err != nil {
		return nil, err
	}

	out, err := JSON6902FromYAML(p, target)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// PatchesPatcher applies JSON6902 or StrategicMergePatch patches into target and
// returns it. It also returns an error, if any.
func PatchesPatcher(patches []string, target []byte) ([]byte, error) {
	var (
		contents    []byte
		err         error
		substituted []string
	)

	for _, patchString := range patches {
		if strings.HasPrefix(patchString, "@") {
			filename := patchString[1:]

			// Try to decrypt patch with sops first.
			contents, err = decrypt.DecryptYamlWithSops(filename)
			if err != nil {
				// If it fails, read the file as is.
				contents, err = os.ReadFile(filename)
				if err != nil {
					return nil, err
				}
			}

			p, err := substitute.SubstituteEnvFromByte(contents)
			if err != nil {
				return nil, err
			}

			substituted = append(substituted, string(p))
		} else {
			substituted = append(substituted, patchString)
		}
	}

	parsedPatches, err := configpatcher.LoadPatches(substituted)
	if err != nil {
		return nil, err
	}

	output, err := configpatcher.Apply(configpatcher.WithBytes(target), parsedPatches)
	if err != nil {
		return nil, err
	}

	cfg, err := output.Bytes()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
