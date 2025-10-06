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
	data := map[string]*tconfig.VersionContract{
		"node1": {Major: 1, Minor: 11},
		"node2": {Major: 1, Minor: 11},
		"node3": {Major: 1, Minor: 12},
	}

	for node, version := range data {
		basecfg, err := os.ReadFile("testdata/" + node + "_basecfg.yaml")
		if err != nil {
			t.Fatal(err)
		}
		in, err := os.ReadFile("testdata/" + node + "_input.yaml")
		if err != nil {
			t.Fatal(err)
		}
		expected, err := os.ReadFile("testdata/" + node + "_expected.yaml")
		if err != nil {
			t.Fatal(err)
		}

		var n config.Node
		if err := yaml.Unmarshal(in, &n); err != nil {
			t.Fatal(err)
		}

		res, err := AddMultiDocs(&n, "metal", basecfg, version)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(res, expected) {
			t.Error(diff.Diff(string(res), string(expected)))
		}
	}
}
