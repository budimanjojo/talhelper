package substitute

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"
)

func SubstituteRelativePaths(configFilePath string, yamlContent []byte) ([]byte, error) {
	// Get the directory of the YAML file
	yamlDir := filepath.Dir(configFilePath)

	// Create a scanner to read through the YAML content
	scanner := bufio.NewScanner(bytes.NewReader(yamlContent))

	// Buffer to hold the processed lines
	var processedLines []string

	// Start reading the file, line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Look for lines containing "@"
		if strings.Contains(line, "@") {
			// Split by "@" to isolate the relative path
			parts := strings.SplitN(line, "@", 2)

			if len(parts) == 2 && len(strings.TrimSpace(parts[1])) > 0 {
				// Get the relative path and resolve to absolute
				relativePath := strings.TrimSpace(parts[1])
				absolutePath := filepath.Join(yamlDir, relativePath)

				// Reconstruct the line with the absolute path
				line = parts[0] + "@" + absolutePath
			}
		}

		// Append the processed line to the result
		processedLines = append(processedLines, line)
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Join all processed lines back into a single string and convert to []byte
	result := strings.Join(processedLines, "\n")

	// Return the processed file
	return []byte(result), nil
}
