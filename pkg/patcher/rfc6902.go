package patcher

import (
	jsonpatch "github.com/evanphx/json-patch"
	yamljson "sigs.k8s.io/yaml"
)

// JSON6902FromYAML applies JSON6902 patch into target yaml bytes.
// It also returns an error, if any.
func JSON6902FromYAML(patch, target []byte) ([]byte, error) {
	return json6902(patch, target, true, true, true)
}

// json6902 applies JSON6902 patch into target bytes.
// Supports `yaml` or `json` for `target`, `patch`, and returned bytes.
// It also returns an error, if any.
func json6902(patch, target []byte, isPatchYAML, isTargetYAML, returnYAML bool) ([]byte, error) {
	var err error
	if isPatchYAML {
		patch, err = yamljson.YAMLToJSON(patch)
		if err != nil {
			return nil, err
		}
	}

	if isTargetYAML {
		target, err = yamljson.YAMLToJSON(target)
		if err != nil {
			return nil, err
		}
	}

	decoded, err := jsonpatch.DecodePatch(patch)
	if err != nil {
		return nil, err
	}

	json, err := decoded.Apply(target)
	if err != nil {
		return nil, err
	}

	if returnYAML {
		json, err = yamljson.JSONToYAML(json)
		if err != nil {
			return nil, err
		}
	}
	return json, nil
}
