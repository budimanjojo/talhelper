package patcher

import (
	"os"
	"reflect"
	"testing"

	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
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
	var m map[any]any
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
	os.Setenv("SOPS_AGE_KEY", "AGE-SECRET-KEY-172FENV3SDP8JSRRX2SWTA9JQMAW7MW3GSKJ2JZDNXS4GVFAS5STQUW8WN4")

	patchList := []string{
		"@testdata/patch.json",
		"@testdata/strategic.yaml",
		"@testdata/patch.yaml",
		"@testdata/encrypted.sops.yaml",
		"@testdata/emptyfile.yaml",
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

func TestPatchesPatherTemplating(t *testing.T) {
	data := []string{"templating0", "templating1", "templating2"}

	for _, d := range data {
		base, err := os.ReadFile("testdata/" + d + "_base.yaml")
		if err != nil {
			t.Fatal(err)
		}

		result, err := PatchesPatcher([]string{"@./testdata/" + d + "_input.yaml"}, base)
		if err != nil {
			t.Fatal(err)
		}

		expected, err := os.ReadFile("testdata/" + d + "_expected.yaml")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("%s\ngot:\n%v\nwant:\n%v", d, string(result), string(expected))
		}
	}
}
