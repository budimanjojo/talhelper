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

	if isSopsEncrypted(m) {
		decrypted, err := decrypt.Data(data, "yaml")
		if err != nil {
			return nil, err
		}
		return decrypted, nil
	}

	return data, nil
}

func isSopsEncrypted(data *sopsFile) bool {
	return len(data.Sops) != 0
}
