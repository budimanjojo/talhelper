package substitute

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// SubstituteFileContent will read and return the content of a file if `value` is string prefixed with `@`
// followed by a path. Otherwise the value will be returned as it is.
// The content will also be envsubst-ed if `envsubst` is `true`.
// It will also returns an error, if any.
func SubstituteFileContent(value string, envsubst bool) (string, error) {
	if strings.HasPrefix(value, "@") {
		slog.Debug(fmt.Sprintf("getting file content of %s", value))
		filename := value[1:]

		contents, err := os.ReadFile(filename)
		if err != nil {
			return "", err
		}

		if envsubst {
			substituted, err := SubstituteEnvFromByte(contents)
			if err != nil {
				return "", err
			}
			return string(substituted), nil
		}

		return string(contents), nil
	}

	return value, nil
}
