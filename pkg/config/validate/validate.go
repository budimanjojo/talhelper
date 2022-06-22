package validate

import (
	"io/ioutil"

	validator "github.com/gookit/validate"
)

func ValidateFromByte(source []byte) (validator.Errors, error){
	return runValidate(source)
}

func ValidateFromFile(path string) (validator.Errors, error) {
	byte, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return runValidate(byte)
}

func runValidate(source []byte) (validator.Errors, error) {
	c, err := NewFromByte(source)
	if err != nil {
		return nil, err
	}
	v := validator.Struct(&c)
	v.StopOnError = false
	if v.Validate() {
		return nil, nil
	} else {
		return v.Errors, err
	}
}
