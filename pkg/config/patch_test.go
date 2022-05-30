package config

import (
	"fmt"
	"reflect"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestMergePatchSlices(t *testing.T) {
	// expected := map{[op:add path:/a/b value:c]} map[op:remove path:/a/b]}]
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
	fmt.Println(expected)
	fmt.Println(string(result))
	if expected != string(result) {
		t.Errorf("got %s, want %s", expected, string(result))

	}
}
