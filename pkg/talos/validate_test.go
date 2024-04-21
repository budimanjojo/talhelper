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

	dataC := []byte(`cluster:
  controlPlane:
    endpoint: https://1.1.1.1:6443
machine:
  type: controlplane
  features:
    hostDNS:
      enabled: true
      forwardKubeDNSToHost: true
`)

	noErr := ValidateConfigFromBytes(data, "metal")
	if noErr == nil {
		t.Errorf("got %s, want %s", noErr, "error")
	}

	err := ValidateConfigFromBytes(dataC, "container")
	if err != nil {
		t.Errorf("got %s, want %s", err, "")
	}
}
