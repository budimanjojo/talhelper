package config

import (
	"testing"
)

func TestGetOutputFileName(t *testing.T) {
	tests := []struct {
		config      *TalhelperConfig
		node        *Node
		expected    string
		shouldError bool
	}{
		{
			config: &TalhelperConfig{ClusterName: "myCluster"},
			node: &Node{
				Hostname:     "node1",
				IPAddress:    "1.1.1.1",
				ControlPlane: true,
				NodeConfigs: NodeConfigs{
					FilenameTmpl: "{{.ClusterName}}-{{.Hostname}}.yaml",
				},
			},
			expected: "myCluster-node1.yaml",
		},
		{
			config: &TalhelperConfig{ClusterName: "myCluster"},
			node: &Node{
				Hostname:     "node1",
				IPAddress:    "1.1.1.1",
				ControlPlane: true,
				NodeConfigs: NodeConfigs{
					FilenameTmpl: "{{.Hostname}}-{{.Role}}.yaml",
				},
			},
			expected: "node1-controlplane.yaml",
		},
		{
			config: &TalhelperConfig{},
			node: &Node{
				Hostname:     "",
				IPAddress:    "1.1.1.1",
				ControlPlane: true,
				NodeConfigs: NodeConfigs{
					FilenameTmpl: "{{.ClusterName}}-{{.Hostname}}.yaml",
				},
			},
			expected: "-.yaml",
		},
		{
			config: &TalhelperConfig{ClusterName: "myCluster"},
			node: &Node{
				Hostname:     "node1",
				IPAddress:    "1.1.1.1",
				ControlPlane: true,
				NodeConfigs: NodeConfigs{
					FilenameTmpl: "{{.Invalid}}.yaml",
				},
			},
			shouldError: true,
		},
	}

	for k, test := range tests {
		result, err := test.node.GetOutputFileName(test.config)
		switch test.shouldError {
		case true:
			if err == nil {
				t.Errorf("tests[%v]\ngot : nil error\nwant: not nil error", k)
			}
		case false:
			if err != nil {
				t.Fatal(err)
			}
		}
		if result != test.expected {
			t.Errorf("tests[%v]\ngot : %v\nwant: %v", k, result, test.expected)
		}
	}
}
