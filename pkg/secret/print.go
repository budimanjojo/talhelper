package secret

import (
	"bytes"
	"fmt"

	"github.com/siderolabs/talos/pkg/machinery/config/generate/secrets"
	"gopkg.in/yaml.v3"
)

// PrintSecretBundle prints the generated `SecretsBundle` into the terminal.
// It returns an error, if any.
func PrintSecretBundle(secret *secrets.Bundle) error {
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	err := encoder.Encode(secret)
	if err != nil {
		return err
	}

	fmt.Print(buf.String())
	return nil
}
