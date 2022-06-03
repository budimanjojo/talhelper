package config

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	data := []byte(`cluster:
  controlPlane:
    endpoint: https://1.1.1.1:6443
machine:
  type: controlplane
`)

	noErr := validateConfig(data, "metal")
	if noErr == nil {
		t.Errorf("got %s, want %s", noErr, "error")
	}

	err := validateConfig(data, "container")
	if err != nil {
		t.Errorf("got %s, want %s", err, "")
	}
}
