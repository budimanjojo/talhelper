package talos

import (
	"github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate/secrets"
)

// NewSecretsBundleFromConfig creates secrets bundle using existing config.
func NewSecretBundleFromCfg(clock secrets.Clock, cfg config.Provider) *secrets.Bundle {
	return secrets.NewBundleFromConfig(clock, cfg)
}

// NewSecretsBundle creates secrets bundle generating all secrets or reading from the input options if provided
func NewSecretBundle(clock secrets.Clock, vc config.VersionContract) (*secrets.Bundle, error) {
	return secrets.NewBundle(clock, &vc)
}
