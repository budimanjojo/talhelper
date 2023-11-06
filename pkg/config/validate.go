package config

import (
	"os"
)

type Warning struct {
	Kind    string
	Field   string
	Message string
}

type Warnings []*Warning

type Error struct {
	Kind    string
	Field   string
	Message error
}

type Errors []*Error

func ValidateFromByte(source []byte) (Errors, Warnings, error) {
	return runValidate(source)
}

func ValidateFromFile(path string) (Errors, Warnings, error) {
	byte, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	return runValidate(byte)
}

// Validate returns `Errors` if the given `TalhelperConfig` is not
// correct
func (c TalhelperConfig) Validate() (Errors, Warnings) {
	var result Errors
	var warns Warnings
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
		checkNodeHostname(node, k, &result)
		checkNodeLabels(node, k, &result)
		checkNodeMachineDisks(node, k, &result)
		checkNodeMachineFiles(node, k, &result)
		checkNodeExtensions(node, k, &result, &warns)
		checkNodeSchematic(node, k, &result)
		checkNodeNameServers(node, k, &result)
		checkNodeNetworkInterfaces(node, k, &result)
		checkNodeConfigPatches(node, k, &result)
	}
	return result, warns
}

func runValidate(source []byte) (Errors, Warnings, error) {
	c, err := NewFromByte(source)
	if err != nil {
		return nil, nil, err
	}
	errors, warnings := c.Validate()
	return errors, warnings, nil
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

func (warns Warnings) HasField(field string) bool {
	for _, warn := range warns {
		if warn.Field == field {
			return true
		}
	}
	return false
}

func (warns *Warnings) Append(warn *Warning) *Warnings {
	*warns = append(*warns, warn)
	return warns
}
