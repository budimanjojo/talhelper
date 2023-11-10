package substitute

import (
	"os"
	"testing"
)

func TestLoadEnvFromFiles(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "AGE-SECRET-KEY-172FENV3SDP8JSRRX2SWTA9JQMAW7MW3GSKJ2JZDNXS4GVFAS5STQUW8WN4")

	files := []string{"testdata/file1.yml", "./testdata/file2.yaml", "./testdata/file3.sops.yaml"}
	expected := map[string]string{
		"env1":          "hello",
		"env2":          "world",
		"env3":          "this is value",
		"enc_hello_env": "hello",
	}
	if err := LoadEnvFromFiles(files); err != nil {
		t.Fatal(err)
	}

	for k, v := range expected {
		if result, _ := os.LookupEnv(k); result != v {
			t.Errorf("%s: got %s, want %s", k, result, v)
		}
	}
}

func TestSubstituteEnvFromYaml(t *testing.T) {
	env := `env1: value1
env2: "true"
env3: 123
default: default value
`

	file := `a: ${env1}
b: "${env2:+$default}" ## commentb
## this is comment
# another comment
c: "${env3-$default}" # commentc
d: ${env4:=$default}
`

	expected := `a: value1
b: "default value"
c: "123"
d: default value
`

	err := LoadEnv([]byte(env))
	if err != nil {
		t.Fatal(err)
	}

	result, _ := SubstituteEnvFromByte([]byte(file))
	if expected != string(result) {
		t.Errorf("got %s, want %s", string(result), expected)
	}
}

func TestLoadEnv(t *testing.T) {
	file := `env1: value1
env2: "true"
env3: 123
default: default value
# env2: "false"
`
	expected := map[string]string{
		"env1":    "value1",
		"env2":    "true",
		"env3":    "123",
		"default": "default value",
	}

	err := LoadEnv([]byte(file))
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range expected {
		if result, _ := os.LookupEnv(k); result != v {
			t.Errorf("%s: got %s, want %s", k, result, v)
		}
	}
}
