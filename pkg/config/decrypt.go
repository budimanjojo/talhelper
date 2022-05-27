package config

import (
	"fmt"
	"os"

	"go.mozilla.org/sops/v3/decrypt"
	"sigs.k8s.io/yaml"
)

type sopsFile struct {
	Sops map[string]interface{} `yaml:"sops"`
}

func DecryptYamlWithSops(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}

	var m sopsFile

	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshalfile: %s", err)
	}

	if isSopsEncrypted(m) {
		decrypted, err := decrypt.Data(data, "yaml")
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt encrypted file: %s", err)
		}
		return decrypted, nil
	}

	return data, nil
}

func isSopsEncrypted(data sopsFile) bool {
	if len(data.Sops) != 0 {
		return true
	}
	return false
}
