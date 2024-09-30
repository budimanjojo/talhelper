package substitute

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// SubstituteRelativePaths will replace all relative paths in the config file to new paths,
// relative to the working dir from which the CLI has been called.
// When using the `--config-file` flag to point to a file not in the current dir,
// relative path evaluation would fail. This function basically prepends the path to the config
// file to the relative paths in the config file so that their evaluation no longer fails.
func SubstituteRelativePaths(configFilePath string, yamlContent []byte) ([]byte, error) {
	// Get the directory of the YAML file
	yamlDir := filepath.Dir(configFilePath)

	// Parse the YAML content
	var data interface{}
	err := yaml.Unmarshal(yamlContent, &data)
	if err != nil {
		return nil, err
	}

	// Process the data
	data = processNode(data, []string{}, yamlDir)

	// Marshal back to YAML
	newYamlContent, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return newYamlContent, nil
}

func processNode(node interface{}, path []string, yamlDir string) interface{} {
	switch n := node.(type) {
	case map[interface{}]interface{}:
		newMap := make(map[interface{}]interface{})
		for k, v := range n {
			keyStr := fmt.Sprintf("%v", k)
			newPath := append(path, keyStr)
			newMap[k] = processNode(v, newPath, yamlDir)
		}
		return newMap

	case []interface{}:
		newArray := make([]interface{}, len(n))
		for i, v := range n {
			newPath := append(path, fmt.Sprintf("[%d]", i))
			newArray[i] = processNode(v, newPath, yamlDir)
		}
		return newArray

	case string:
		if shouldSubstitute(path) {
			if strings.HasPrefix(n, "@") {
				parts := strings.SplitN(n, "@", 2)
				if len(parts) == 2 && len(strings.TrimSpace(parts[1])) > 0 {
					relativePath := strings.TrimSpace(parts[1])
					absolutePath := filepath.Join(yamlDir, relativePath)
					return parts[0] + "@" + absolutePath
				}
			}
		}
		return n

	default:
		return n
	}
}

func shouldSubstitute(path []string) bool {
	for _, p := range path {
		if p == "machineFiles" || p == "patches" || p == "extraManifests" {
			return true
		}
	}
	return false
}
