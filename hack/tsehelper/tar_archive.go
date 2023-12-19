// tar_archive.go
package main

import (
	"archive/tar"
	"bytes"
	"io"
	"strings"
)

// processTarArchive processes a tar archive and returns a slice of strings or an error.
func processTarArchive(tarData []byte) ([]string, error) {
	// Create a new reader from the tar data
	tarReader := tar.NewReader(bytes.NewReader(tarData))

	// Iterate through the files in the tar archive
	var exts []string
	for {
		_, err := tarReader.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		if err != nil {
			return nil, err
		}
		var fileContent bytes.Buffer
		_, err = io.Copy(&fileContent, tarReader)
		if err != nil {
			return nil, err
		}
		// Process the file content, split by new line
		lines := strings.Split(fileContent.String(), "\n")

		for _, line := range lines {
			// Skip empty lines
			if line == "" {
				continue
			}
			exts = append(exts, line)
		}
	}

	return exts, nil
}
