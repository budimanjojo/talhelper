package substitute

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSubstituteRelativePaths(t *testing.T) {
	// Define the path to the YAML file (for testing purposes, we use a dummy path)
	configFilePath := "/path/to/config.yaml"

	// Define test cases
	tests := []struct {
		name           string
		yamlContent    string
		expectedOutput string
	}{
		{
			name: "Substitute in machineFiles",
			yamlContent: `
machineFiles:
  - content: "@./relative/path/file.txt"
    permissions: 0o644
    path: /var/etc/tailscale/auth.env
    op: create
`,
			expectedOutput: `
machineFiles:
  - content: "@/path/to/relative/path/file.txt"
    permissions: 0o644
    path: /var/etc/tailscale/auth.env
    op: create
`,
		},
		{
			name: "Substitute in patches",
			yamlContent: `
patches:
  - "@./relative/patch.yaml"
`,
			expectedOutput: `
patches:
  - "@/path/to/relative/patch.yaml"
`,
		},
		{
			name: "No substitution outside specified sections",
			yamlContent: `
nodes:
  - hostname: kworker1
    ipAddress: "@./should/not/replace"
`,
			expectedOutput: `
nodes:
  - hostname: kworker1
    ipAddress: "@./should/not/replace"
`,
		},
		{
			name: "Substitute in nested machineFiles",
			yamlContent: `
nodes:
  - hostname: kmaster1
    machineFiles:
      - content: "@./nested/path/file.conf"
        path: /etc/config.conf
`,
			expectedOutput: `
nodes:
  - hostname: kmaster1
    machineFiles:
      - content: "@/path/to/nested/path/file.conf"
        path: /etc/config.conf
`,
		},
		{
			name: "Substitute in nested patches",
			yamlContent: `
controlPlane:
  patches:
    - "@./control/plane/patch.yaml"
`,
			expectedOutput: `
controlPlane:
  patches:
    - "@/path/to/control/plane/patch.yaml"
`,
		},
		{
			name: "Multiple substitutions",
			yamlContent: `
machineFiles:
  - content: "@./file1.txt"
patches:
  - "@./patch1.yaml"
nodes:
  - hostname: node1
    ipAddress: "@./should/not/replace"
    machineFiles:
      - content: "@./node1/file2.txt"
    patches:
      - "@./node1/patch2.yaml"
`,
			expectedOutput: `
machineFiles:
  - content: "@/path/to/file1.txt"
patches:
  - "@/path/to/patch1.yaml"
nodes:
  - hostname: node1
    ipAddress: "@./should/not/replace"
    machineFiles:
      - content: "@/path/to/node1/file2.txt"
    patches:
      - "@/path/to/node1/patch2.yaml"
`,
		},
		{
			name: "No substitution when '@' is not at the beginning",
			yamlContent: `
machineFiles:
  - content: "example@./should/not/replace"
`,
			expectedOutput: `
machineFiles:
  - content: "example@./should/not/replace"
`,
		},
		{
			name: "Handle empty '@' reference",
			yamlContent: `
machineFiles:
  - content: "@"
`,
			expectedOutput: `
machineFiles:
  - content: "@"
`,
		},
		{
			name: "Handle missing relative path after '@'",
			yamlContent: `
machineFiles:
  - content: "@ "
`,
			expectedOutput: `
machineFiles:
  - content: "@ "
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputBytes, err := SubstituteRelativePaths(configFilePath, []byte(tt.yamlContent))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			var expectedData interface{}
			err = yaml.Unmarshal([]byte(tt.expectedOutput), &expectedData)
			if err != nil {
				t.Fatalf("Failed to unmarshal expected output: %v", err)
			}

			var actualData interface{}
			err = yaml.Unmarshal(outputBytes, &actualData)
			if err != nil {
				t.Fatalf("Failed to unmarshal actual output: %v", err)
			}

			if !reflect.DeepEqual(actualData, expectedData) {
				// For better error messages, re-marshal the data to YAML strings
				expectedYAML, _ := yaml.Marshal(expectedData)
				actualYAML, _ := yaml.Marshal(actualData)
				t.Errorf("Output mismatch.\nExpected:\n%s\nGot:\n%s", expectedYAML, actualYAML)
			}
		})
	}
}
