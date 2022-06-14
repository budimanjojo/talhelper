package decrypt

import (
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
		return nil, err
	}

	var m *sopsFile

	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	if m.isEncrypted() {
		decrypted, err := decrypt.Data(data, "yaml")
		if err != nil {
			return nil, err
		}
		return decrypted, nil
	}

	return data, nil
}

func (s *sopsFile) isEncrypted() bool {
	return len(s.Sops) != 0
}
