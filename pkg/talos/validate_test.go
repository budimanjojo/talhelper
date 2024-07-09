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
  ca:
    crt: hehe
    key: hehe
  features:
    hostDNS:
      enabled: true
      forwardKubeDNSToHost: true
`)

	err := ValidateConfigFromBytes(data, "metal")
	if err == nil {
		t.Errorf("got %s, want %s", err, "error")
	}

	noErr := ValidateConfigFromBytes(dataC, "container")
	if noErr != nil {
		t.Errorf("got %s, want %s", noErr, "")
	}
}
