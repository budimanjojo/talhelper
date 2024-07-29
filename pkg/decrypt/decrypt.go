package decrypt

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/getsops/sops/v3/decrypt"
	"sigs.k8s.io/yaml"
)

type sopsFile struct {
	Sops map[string]interface{} `yaml:"sops"`
}

// DecryptYamlWithSops reads a `sops` encrypted `yaml` file path
// and decrypt the content using `sops/v3/decrypt.Data`.
// The unencrypted data will be returned bytes.
// Data will be returned as it is if file is not encrypted with
// `sops`. Error will be returned when decryption fails.
func DecryptYamlWithSops(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// return the data immediately if the file is empty
	if len(data) == 0 {
		return data, nil
	}

	var m *sopsFile

	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	if m.isEncrypted() {
		slog.Debug(fmt.Sprintf("%s is SOPS encrypted, decrypting", filePath))
		decrypted, err := decrypt.Data(data, "yaml")
		if err != nil {
			return nil, fmt.Errorf("SOPS decryption failed: %w", err)
		}
		return decrypted, nil
	}

	return data, nil
}

// isEncrypted returns true if `sops` key exists.
func (s *sopsFile) isEncrypted() bool {
	return len(s.Sops) != 0
}
