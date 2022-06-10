package config

import (
	"os"

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

func SubstituteEnvFromYaml(file []byte) ([]byte, error) {
	data, err := envsubst.BytesRestricted(file, true, true)
	if err != nil {
		return nil, err
	}

	return data, nil
}
