package config

import (
	"os"
	"testing"
)

func TestSubstituteEnvFromYaml(t *testing.T) {
	env := `env1: value1
env2: "true"
env3: 123
default: default value
`

	file := `a: ${env1}
b: "${env2:+$default}"
c: "${env3-$default}"
d: ${env4:=$default}
`

	expected := `a: value1
b: "default value"
c: "123"
d: default value
`

	result, _ := SubstituteEnvFromYaml([]byte(env), []byte(file))
	if expected != string(result) {
		t.Errorf("got %s, want %s", string(result), expected)

	}
}

func TestLoadEnv(t *testing.T) {
	file := `env1: value1
env2: "true"
env3: 123
default: default value
`
	expected := map[string]string{
		"env1": "value1",
		"env2": "true",
		"env3": "123",
		"default": "default value",
	}

	loadEnv([]byte(file))
	for k, v := range expected {
		if result, _ := os.LookupEnv(k); result != v {
			t.Errorf("got %s, want %s", result, v)
		}
	}

}
