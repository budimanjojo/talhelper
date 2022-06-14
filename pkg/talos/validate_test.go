package talos

import (
	"testing"
)

func TestValidateConfigFromBytes(t *testing.T) {
	data := []byte(`cluster:
  controlPlane:
    endpoint: https://1.1.1.1:6443
machine:
  type: controlplane
`)

	noErr := ValidateConfigFromBytes(data, "metal")
	if noErr == nil {
		t.Errorf("got %s, want %s", noErr, "error")
	}

	err := ValidateConfigFromBytes(data, "container")
	if err != nil {
		t.Errorf("got %s, want %s", err, "")
	}
}
