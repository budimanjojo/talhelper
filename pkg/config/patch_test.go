package config

import (
	"reflect"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestMergePatchSlices(t *testing.T) {
	var test []map[string]interface{}
	map1 := map[string]interface{}{"op": "add", "path": "/a/b", "value": "c"}
	map2 := map[string]interface{}{"op": "remove", "path": "/a/b"}
	expected := append(test, map1, map2)

	data := `controlPlane:
  patches:
    - op: add
      path: /a/b
      value: c
  encryptedPatches:
    - op: remove
      path: /a/b`
	var m TalhelperConfig
	yaml.Unmarshal([]byte(data), &m)
	result := mergePatchSlices(m.ControlPlane.Patches, m.ControlPlane.EncryptedPatches)
	if len(result) != len(expected) {
		t.Errorf("got %d, want %d", len(result), len(expected))
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %t, want true", reflect.DeepEqual(result, expected))
	}

}

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

	result, _ := applyPatchFromYaml([]byte(patch), []byte(file))
	if expected != string(result) {
		t.Errorf("got %s, want %s", string(result), expected)

	}
}

func TestApplyInlinePatchFromYaml(t *testing.T) {
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

	result, _ := applyInlinePatchFromYaml([]byte(patch), []byte(file))
	if expected != string(result) {
		t.Errorf("got %s, want %s", string(result), expected)

	}
}
