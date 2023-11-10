package substitute

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/a8m/envsubst"
	"github.com/budimanjojo/talhelper/pkg/decrypt"
	"github.com/joho/godotenv"
)

// LoadEnvFromFiles read yaml data from list of filepaths and sets
// environment variable named by the key. It will try to decrypt
// with `sops` if the file is encrypted and skips if file doesn't
// exist. It returns an error, if any.
func LoadEnvFromFiles(files []string) error {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			env, err := decrypt.DecryptYamlWithSops(file)
			if err != nil {
				return fmt.Errorf("trying to decrypt %s with sops: %s", file, err)
			}

			if err := LoadEnv(env); err != nil {
				return fmt.Errorf("trying to load env from %s: %s", file, err)
			}
		} else if errors.Is(err, os.ErrNotExist) {
			continue
		} else {
			return fmt.Errorf("trying to stat %s: %s", file, err)
		}
	}
	return nil
}

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
	stripped := re.ReplaceAllFunc(file, func(b []byte) []byte {
		re := regexp.MustCompile("^['\"].+['\"]|^[a-zA-Z0-9]")
		if re.Match(b) {
			return b
		} else {
			return []byte("\n")
		}
	})

	var final bytes.Buffer
	for _, line := range bytes.Split(stripped, []byte("\n")) {
		if len(bytes.TrimSpace(line)) > 0 {
			final.WriteString(string(line) + "\n")
		}
	}

	return final.Bytes()
}
