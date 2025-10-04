package talos

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/kylelemons/godebug/diff"
	"gopkg.in/yaml.v3"
)

func TestAddMultiDocs(t *testing.T) {
	for _, node := range []string{"node1", "node2"} {
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

		res, err := AddMultiDocs(&n, "metal", basecfg)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(res, expected) {
			fmt.Println(diff.Diff(string(res), string(expected)))
			t.Errorf("\ngot:\n%s\n\nwant:\n%s", res, expected)
		}
	}
}
