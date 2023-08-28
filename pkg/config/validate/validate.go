package validate

import (
	"os"

	validator "github.com/gookit/validate"
)

// ValidateFromByte reads yaml bytes and validates the data.
// It returns all the incorrect values and an error, if any.
func ValidateFromByte(source []byte) (validator.Errors, error){
	return runValidate(source)
}

// ValidateFromFile reads yaml file path and validates the data.
// It returns all the incorrect values and an error, if any.
func ValidateFromFile(path string) (validator.Errors, error) {
	byte, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return runValidate(byte)
}

// runValidate reads yaml bytes and validates a `struct` of Config
// using `gookit/validate`. It returns all the incorrect values an
// and error, if any.
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
