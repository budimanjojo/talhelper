package talos

import (
	"github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
)

func NewSecretFromCfg(clock generate.Clock, cfg config.Provider) *generate.SecretsBundle {
	return generate.NewSecretsBundleFromConfig(clock, cfg)
}

func NewSecretBundle(clock generate.Clock, opts ...generate.GenOption) (*generate.SecretsBundle, error) {
	return generate.NewSecretsBundle(clock, opts...)
}
