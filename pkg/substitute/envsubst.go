package substitute

import (
	"os"
	"regexp"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
)

// LoadEnv reads yaml data and sets environment variable named
// by the key. It returns an error, if any.
func LoadEnv(file []byte) error {
	mFile, err := godotenv.Unmarshal(string(file))
	if err != nil {
		return nil
	}

	for k, v := range mFile {
		os.Setenv(k, v)
	}
	return nil
}

// SubstituteEnvFromByte reads yaml bytes and do `envsubst` on
// them. The substituted bytes will be returned. It returns an
// error, if any.
func SubstituteEnvFromByte(file []byte) ([]byte, error) {
	filtered := stripYamlComment(file)
	data, err := envsubst.BytesRestricted(filtered, true, true)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// stripYamlComment takes yaml bytes and returns them back with
// comments stripped.
func stripYamlComment(file []byte) []byte {
	// FIXME use better logic than regex.
	re := regexp.MustCompile(".?#.*\n")
	first := re.ReplaceAllFunc(file, func(b []byte) []byte {
		re := regexp.MustCompile("^['\"].+['\"]|^[a-zA-Z1-9]")
		if re.Match(b) {
			return b
		} else {
			return []byte("\n")
		}
	})

	return first
}
