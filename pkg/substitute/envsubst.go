package substitute

import (
	"os"
	"regexp"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
)

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

func SubstituteEnvFromByte(file []byte) ([]byte, error) {
	filtered := stripYamlComment(file)
	data, err := envsubst.BytesRestricted(filtered, true, true)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func stripYamlComment(file []byte) ([]byte) {
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
