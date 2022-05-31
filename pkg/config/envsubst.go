package config

import (
	"os"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
)

func loadEnv(file []byte) error {
	mFile, err := godotenv.Unmarshal(string(file))
	if err != nil {
		return nil
	}
	for k, v := range mFile {
		os.Setenv(k, v)
	}
	return nil
}

func SubstituteEnvFromYaml(env, file []byte) ([]byte, error) {
	err := loadEnv(env)
	data, err := envsubst.BytesRestricted(file, true, true)
	if err != nil {
		return nil, err
	}
	return data, nil
}
