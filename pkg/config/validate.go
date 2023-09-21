package config

import (
	"os"
)

type Error struct {
	Kind    string
	Field   string
	Message error
}

type Errors []*Error

func ValidateFromByte(source []byte) (Errors, error) {
	return runValidate(source)
}

func ValidateFromFile(path string) (Errors, error) {
	byte, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return runValidate(byte)
}

// Validate returns `Errors` if the given `TalhelperConfig` is not
// correct
func (c TalhelperConfig) Validate() Errors {
	var result Errors
	checkRequiredCfg(c, &result)
	checkSupportedTalosVersion(c, &result)
	checkSupportedK8sVersion(c, &result)
	checkTalosEndpoint(c, &result)
	checkDomain(c, &result)
	checkClusterNets(c, &result)
	checkCNIConfig(c, &result)
	checkControlPlane(c, &result)
	checkWorker(c, &result)
	for k, node := range c.Nodes {
		checkNodeRequiredCfg(node, k, &result)
		checkNodeIPAddress(node, k, &result)
		checkNodeLabels(node, k, &result)
		checkNodeMachineDisks(node, k, &result)
		checkNodeMachineFiles(node, k, &result)
		checkNodeExtensions(node, k, &result)
		checkNodeNameServers(node, k, &result)
		checkNodeNetworkInterfaces(node, k, &result)
		checkNodeConfigPatches(node, k, &result)
	}
	return result
}

func runValidate(source []byte) (Errors, error) {
	c, err := NewFromByte(source)
	if err != nil {
		return nil, err
	}
	errors := c.Validate()
	return errors, nil
}

func (errs Errors) HasField(field string) bool {
	for _, err := range errs {
		if err.Field == field {
			return true
		}
	}
	return false
}

func (errs *Errors) Append(err *Error) *Errors {
	*errs = append(*errs, err)
	return errs
}
