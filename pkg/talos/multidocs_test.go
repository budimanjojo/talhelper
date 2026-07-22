package talos

import (
	"bytes"
	"os"
	"testing"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/kylelemons/godebug/diff"
	tconfig "github.com/siderolabs/talos/pkg/machinery/config"
	"gopkg.in/yaml.v3"
)

func TestAddMultiDocs(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		version *tconfig.VersionContract
	}{
		{"ingress firewall without multidoc", "node1", &tconfig.VersionContract{Major: 1, Minor: 11}},
		{"user volumes without multidoc", "node2", &tconfig.VersionContract{Major: 1, Minor: 11}},
		{"ingress firewall with multidoc", "node3", &tconfig.VersionContract{Major: 1, Minor: 12}},
		{"no multidocs", "node4", &tconfig.VersionContract{Major: 1, Minor: 11}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basecfg, n, expected := loadMultiDocsTestdata(t, tt.file)

			res, err := AddMultiDocs(&n, "metal", basecfg, tt.version)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(res, expected) {
				t.Error(diff.Diff(string(res), string(expected)))
			}
		})
	}
}

func loadMultiDocsTestdata(t *testing.T, node string) (basecfg []byte, n config.Node, expected []byte) {
	t.Helper()

	basecfg, err := os.ReadFile("testdata/" + node + "_basecfg.yaml")
	if err != nil {
		t.Fatal(err)
	}
	in, err := os.ReadFile("testdata/" + node + "_input.yaml")
	if err != nil {
		t.Fatal(err)
	}
	expected, err = os.ReadFile("testdata/" + node + "_expected.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if err := yaml.Unmarshal(in, &n); err != nil {
		t.Fatal(err)
	}

	return basecfg, n, expected
}
