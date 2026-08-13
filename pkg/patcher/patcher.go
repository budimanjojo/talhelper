package patcher

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"unicode"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/decrypt"
	"github.com/budimanjojo/talhelper/v3/pkg/substitute"
	"github.com/budimanjojo/talhelper/v3/pkg/templating"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	"github.com/siderolabs/talos/pkg/machinery/config/configpatcher"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	"gopkg.in/yaml.v3"
)

// TemplateData is the data passed to templates rendered from external patch
// and manifest files. It embeds the rendered Talos `*v1alpha1.Config` so
// existing templates can keep referencing `.MachineConfig`/`.ClusterConfig`
// via field promotion, and adds `.TalConfig` which exposes the source
// `talconfig.yaml` (e.g. `.TalConfig.Nodes` to reach every node's IP address).
type TemplateData struct {
	*v1alpha1.Config
	TalConfig *config.TalhelperConfig
}

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
func PatchesPatcher(patches []string, target []byte, talconfig *config.TalhelperConfig) ([]byte, error) {
	var (
		contents    []byte
		err         error
		substituted []string
	)

	provider, err := configloader.NewFromBytes(target)
	if err != nil {
		return nil, err
	}
	templateData := TemplateData{Config: provider.RawV1Alpha1(), TalConfig: talconfig}

	for _, patchString := range patches {
		if strings.HasPrefix(patchString, "@") {
			filename := patchString[1:]

			// skip empty file
			empty, err := isEmptyFile(filename)
			if err != nil {
				return nil, err
			}
			if empty {
				slog.Debug(fmt.Sprintf("%s is an empty file, skip applying this patch", filename))
				continue
			}

			// Try to decrypt patch with sops first.
			contents, err = decrypt.DecryptYamlWithSops(filename)
			if err != nil {
				// If it fails, read the file as is.
				contents, err = os.ReadFile(filename)
				if err != nil {
					return nil, err
				}
			}

			// templating first before substitution so it doesn't breaks templating with variables
			// like {{ $var }}. And it will only work for patches in a file too because substitution is
			// being done in config file first, there's nothing I can do about it
			p, err := templating.RenderTemplate[[]byte](string(contents), templateData)
			if err != nil {
				return nil, err
			}

			p, err = substitute.SubstituteEnvFromByte(p)
			if err != nil {
				return nil, err
			}

			substituted = append(substituted, string(p))
		} else {
			patchString, err = templating.RenderTemplate[string](patchString, templateData)
			if err != nil {
				return nil, err
			}

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

// YamlBytesPatcher applies StrategicMergePatch patches into target and returns
// it. It doesn't do any substitution so it's more efficient. The patch should
// be YAML encoded and can be multi-document machineconfig.
// It also returns an error, if any.
func YamlBytesPatcher(patch, target []byte) ([]byte, error) {
	slog.Debug("patching multidocs configurations to main configurations")
	var patches []configpatcher.Patch

	p, err := configpatcher.LoadPatch(patch)
	if err != nil {
		return nil, err
	}
	patches = append(patches, p)
	output, err := configpatcher.Apply(configpatcher.WithBytes(target), patches)
	if err != nil {
		return nil, err
	}

	return output.Bytes()
}

// isEmptyFile checks if the file is empty or contains only whitespaces.
// It also returns an error, if any.
func isEmptyFile(file string) (bool, error) {
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimFunc(line, unicode.IsSpace) != "" {
			return false, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, nil
	}

	return true, nil
}
