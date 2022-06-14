package patcher

import (
	"testing"

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
	err := yaml.Unmarshal([]byte(patch), &m)
	if err != nil {
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
	err := yaml.Unmarshal([]byte(patch), &m)
	if err != nil {
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
