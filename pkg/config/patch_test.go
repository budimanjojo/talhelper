package config

import (
	"testing"
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
