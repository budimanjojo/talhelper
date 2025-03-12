package substitute

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SubstituteRelativePaths replaces value of special keys from relative paths to absolute paths.
// The substituted paths are relative to the directory of `configFilePath`.
// When using the `--config-file` flag to point to a file not in the current dir,
// relative path evaluation would fail. This function basically prepends the path to the config
// file to the relative paths in the config file so that their evaluation no longer fails.
// It returns an error, if any.
func SubstituteRelativePaths(configFilePath string, yamlContent []byte) ([]byte, error) {
	// Get the directory of the YAML file
	absolutePath, err := filepath.Abs(filepath.Dir(configFilePath))
	if err != nil {
		return nil, err
	}

	// Parse the YAML content
	var data interface{}
	err = yaml.Unmarshal(yamlContent, &data)
	if err != nil {
		return nil, err
	}

	// Process the data
	data = processNode(data, []string{}, absolutePath)

	// Marshal back to YAML
	newYamlContent, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return newYamlContent, nil
}

func processNode(node interface{}, path []string, yamlDir string) interface{} {
	switch n := node.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range n {
			newPath := append(path, k)
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
		should, special := shouldSubstitute(path)
		if should {
			return handleSubstitution(n, yamlDir, special)
		} else {
			return n
		}
	default:
		return n
	}
}

func shouldSubstitute(path []string) (should, special bool) {
	for _, p := range path {
		// this is special case where the key was introduced without needing
		// "@" prefix, instead of breaking changes, we now internally add the
		// prefix instead
		if p == "extraManifests" {
			return true, true
		} else if p == "machineFiles" || p == "patches" || p == "inlineManifests" {
			return true, false
		}
	}
	return false, false
}

func handleSubstitution(val, yamlDir string, special bool) string {
	// we add "@" to the value of special case like "extraManifests" so we can
	// handle them uniformly with all other keys
	if special && !strings.HasPrefix(val, "@") {
		val = "@" + val
	}

	path, found := strings.CutPrefix(val, "@")
	if found {
		path = strings.TrimSpace(path)
		// the value is unchanged if there's nothing after @
		if path == "" {
			return val
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(yamlDir, path)
		}
		return "@" + path
	}

	return val
}
