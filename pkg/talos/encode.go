package talos

import (
	"bytes"

	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	"gopkg.in/yaml.v3"
)

// ReEncodeTalosConfig takes `yaml` encoded bytes and re-encodes it back to
// 2 indentation `yaml` bytes. It also returns an error, if any.
func ReEncodeTalosConfig(f []byte) ([]byte, error) {
	var c *v1alpha1.Config
	err := yaml.Unmarshal(f, &c)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	err = encoder.Encode(c)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
