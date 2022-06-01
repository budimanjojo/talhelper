package config

import (
	jsonpatch "github.com/evanphx/json-patch"
	yamljson "sigs.k8s.io/yaml"
)

func ApplyInlinePatchFromYaml(patch, yaml []byte) (output []byte, err error) {
	jsonPatch, err := yamljson.YAMLToJSON(patch)
	if err != nil {
		return nil, err
	}

	jsonFile, err := yamljson.YAMLToJSON(yaml)
	if err != nil {
		return nil, err
	}

	finalJson, err := jsonpatch.MergePatch(jsonFile, jsonPatch)
	if err != nil {
		return nil, err
	}

	finalYaml, err := yamljson.JSONToYAML(finalJson)
	if err != nil {
		return nil, err
	}
	
	return finalYaml, nil
}

func applyPatchFromYaml(patch, yaml []byte) (output []byte, err error) {
	jsonPatch, err := yamljson.YAMLToJSON(patch)
	if err != nil {
		return nil, err
	}

	jsonFile, err := yamljson.YAMLToJSON(yaml)
	if err != nil {
		return nil, err
	}

	decodedPatch, err := jsonpatch.DecodePatch(jsonPatch)
	if err != nil {
		return nil, err
	}

	finalJson, err := decodedPatch.Apply(jsonFile)
	if err != nil {
		return nil, err
	}

	finalYaml, err := yamljson.JSONToYAML(finalJson)
	if err != nil {
		return nil, err
	}
	
	return finalYaml, nil
}
