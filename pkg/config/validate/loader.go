package validate

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

func NewFromByte(source []byte) (Config, error) {
	return newConfig(source)
}

func NewFromFile(path string) (c Config, err error) {
	source, err := fromFile(path)
	if err != nil {
		return c, err
	}

	return newConfig(source)
}

func fromFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func newConfig(source []byte) (c Config, err error) {
	err = yaml.Unmarshal(source, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
