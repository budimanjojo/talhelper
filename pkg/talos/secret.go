package talos

import (
	"github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1/generate"
)

// NewSecretsBundleFromConfig creates secrets bundle using existing config.
func NewSecretBundleFromCfg(clock generate.Clock, cfg config.Provider) *generate.SecretsBundle {
	return generate.NewSecretsBundleFromConfig(clock, cfg)
}

// NewSecretsBundle creates secrets bundle generating all secrets or reading from the input options if provided
func NewSecretBundle(clock generate.Clock, opts ...generate.GenOption) (*generate.SecretsBundle, error) {
	return generate.NewSecretsBundle(clock, opts...)
}
