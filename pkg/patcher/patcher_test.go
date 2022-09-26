package patcher

import (
	"os"
	"testing"

	"github.com/talos-systems/talos/pkg/machinery/config/configloader"
	"gopkg.in/yaml.v3"
)

func TestApplyPatchFromYaml(t *testing.T) {
	patch := `
- op: add
  path: /a/b
  value: added`

	file := `a:
  b: original
`

	expected := `a:
  b: added
`
	var m []map[interface{}]interface{}
	if err := yaml.Unmarshal([]byte(patch), &m); err != nil {
		t.Fatal(err)
	}

	result, err := YAMLPatcher(m, []byte(file))
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(result) {
		t.Errorf("got %s, want %s", string(result), expected)

	}
}

func TestYAMLInlinePatcher(t *testing.T) {
	patch := `a:
  b:
    c: added
    d: added
`

	file := `a:
  b:
    c: original
  c: original
`

	expected := `a:
  b:
    c: added
    d: added
  c: original
`
	var m map[interface{}]interface{}
	if err := yaml.Unmarshal([]byte(patch), &m); err != nil {
		t.Fatal(err)
	}

	result, err := YAMLInlinePatcher(m, []byte(file))
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(result) {
		t.Errorf("got %s, want %s", string(result), expected)

	}
}

func TestPatchesPatcher(t *testing.T) {
	os.Setenv("foodomain", "foo.com")
	os.Setenv("foodotbar", "foo.bar")

	patchList := []string{
		"@testdata/patch.json",
		"@testdata/strategic.yaml",
		"@testdata/patch.yaml",
		`[{"op":"add","path":"/machine/network/interfaces/0/dhcp","value": false}]`,
	}

	file := []byte(`version: v1alpha1
machine:
  certSANs:
    - ""
  network:
    interfaces:
      - interface: eth0
        dhcp: false
`)

	expected, err := configloader.NewFromBytes([]byte(`cluster: null
machine:
  certSANs:
  - foo.com
  network:
    hostname: foo.bar
    interfaces:
    - addresses:
      - 10.1.2.3/24
      dhcp: false
      interface: eth0
      mtu: 0
  token: ""
  type: ""
version: v1alpha1
`))
	if err != nil {
		t.Fatal(err)
	}

	result, err := PatchesPatcher(patchList, []byte(file))
	if err != nil {
		t.Fatal(err)
	}

	expectedByte, err := expected.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	if string(expectedByte) != string(result) {
		t.Errorf("got %s, want %s", string(result), string(expectedByte))

	}
}
