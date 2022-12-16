package talos

import (
	"bytes"

	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	"gopkg.in/yaml.v3"
)

func ReEncodeTalosConfig(f []byte, c *v1alpha1.Config) ([]byte, error) {
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
