package validate

import (
	"os"

	"github.com/ghodss/yaml"
)

// NewFromByte takes bytes and convert it into Talhelper config.
// It also returns an error, if any.
func NewFromByte(source []byte) (Config, error) {
	return newConfig(source)
}

// NewFromFile takes a file path and convert the contents into Talhelper config.
// It also returns an error, if any.
func NewFromFile(path string) (c Config, err error) {
	source, err := fromFile(path)
	if err != nil {
		return c, err
	}

	return newConfig(source)
}

// fromFile is a wrapper for `os.ReadFile`.
func fromFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// newConfig takes bytes and convert it into Talhelper config.
// It also returns an error, if any.
func newConfig(source []byte) (c Config, err error) {
	err = yaml.Unmarshal(source, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
