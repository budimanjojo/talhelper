package validate

import (
	"testing"
)

func TestValidateFromByte(t *testing.T) {
	data := []byte(`clusterName: "test-cluster"
talosVersion: vv1.0.6
kubernetesVersion: 1.23.6.5
endpoint: https://1.1.1.1:6443
allowSchedulingOnMasters: Truee
cniConfig:
  name: customo
nodes:
  - hostname: master1
    ipAddress: 1.2.3.4.5
    installDisk: /dev/sda
    nameservers:
      - 8.8.8.8
    networkInterfaces:
      - addresses:
          - 1.2.3.4
        mtu: one
        routes:
          - network: 0.0.0.0/0
            gateway: 1.2.3.4.5.6
    configPatches:
      - op: del
        path: /cluster
`)

	found, err := ValidateFromByte(data)
	if err != nil {
		t.Fatal(err)
	}

	expectedErrors := map[string]bool{
		"ClusterName": false,
		"TalosVersion": true,
		"KubernetesVersion": true,
		"CniConfig.Name": true,
		"AllowSchedulingOnMasters": true,
		"Nodes.0.Hostname": false,
		"Nodes.0.IPAddress": true,
		"Nodes.0.ControlPlane": false,
		"Nodes.0.InstallDisk": false,
		"Nodes.0.Nameservers": false,
		"Nodes.0.NetworkInterfaces.0.Interface": true,
		"Nodes.0.NetworkInterfaces.0.Addresses": true,
		"Nodes.0.NetworkInterfaces.0.Mtu": true,
		"Nodes.0.NetworkInterfaces.0.Routes.0.Gateway": true,
		"Nodes.0.ConfigPatches": true,
	}

	for k, v := range expectedErrors {
		if found.HasField(k) != v {
			t.Errorf("%s: got %t, want %t", k, found.HasField(k), v)
		}
	}
}
