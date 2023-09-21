package config

import (
	"os"

	"github.com/go-yaml/yaml"
)

// NewFromByte takes bytes and convert it into Talhelper config.
// It also returns an error, if any.
func NewFromByte(source []byte) (TalhelperConfig, error) {
	return newConfig(source)
}

// NewFromFile takes a file path and convert the contents into Talhelper config.
// It also returns an error, if any.
func NewFromFile(path string) (c TalhelperConfig, err error) {
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
func newConfig(source []byte) (c TalhelperConfig, err error) {
	err = yaml.Unmarshal(source, &c)
	if err != nil {
		return TalhelperConfig{}, err
	}
	return c, nil
}
