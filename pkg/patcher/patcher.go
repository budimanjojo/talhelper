package patcher

import "gopkg.in/yaml.v3"

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
