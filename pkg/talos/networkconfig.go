package talos

import (
	"bytes"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/go-pointer"
	"github.com/siderolabs/talos/pkg/machinery/config/types/network"
	"github.com/siderolabs/talos/pkg/machinery/nethelpers"
	"gopkg.in/yaml.v3"
)

func GenerateNetworkHostnameConfigBytes(name string, stableHostname bool) ([]byte, error) {
	var result [][]byte

	hostname := GenerateNetworkHostnameConfig(name, stableHostname)
	hostnameBytes, err := marshalYaml(hostname)
	if err != nil {
		return nil, err
	}

	return CombineYamlBytes(append(result, hostnameBytes)), nil
}

func GenerateNetworkHostnameConfig(name string, stableHostname bool) *network.HostnameConfigV1Alpha1 {
	result := network.NewHostnameConfigV1Alpha1()

	if stableHostname {
		result.ConfigAuto = pointer.To(nethelpers.AutoHostnameKindStable)
		return result
	} else {
		result.ConfigHostname = name
		// TODO: this is awkward because the Generate API is handling this by default.
		// On version greater than 1.1 above, stable hostname is enabled by default and it will conflict
		// with us setting the hostname field
		result.ConfigAuto = pointer.To(nethelpers.AutoHostnameKindOff)
		return result
	}
}

func GenerateNetworkConfigBytes(ifCfg *config.IngressFirewall) ([]byte, error) {
	var result [][]byte

	defaultAction := GenerateNodeDefaultActionConfig(ifCfg)
	defaultActionBytes, err := marshalYaml(defaultAction)
	if err != nil {
		return nil, err
	}

	result = append(result, defaultActionBytes)

	rules, err := GenerateNodeRuleConfig(ifCfg)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		ruleBytes, err := marshalYaml(rule)
		if err != nil {
			return nil, err
		}

		result = append(result, ruleBytes)
	}

	return CombineYamlBytes(result), nil
}

func GenerateNodeDefaultActionConfig(ifCfg *config.IngressFirewall) *network.DefaultActionConfigV1Alpha1 {
	result := network.NewDefaultActionConfigV1Alpha1()
	result.Ingress = ifCfg.DefaultAction

	return result
}

func GenerateNodeRuleConfig(ifCfg *config.IngressFirewall) ([]*network.RuleConfigV1Alpha1, error) {
	var result []*network.RuleConfigV1Alpha1

	for _, v := range ifCfg.NetworkRules {
		rule := network.NewRuleConfigV1Alpha1()
		rule.MetaName = v.Name
		rule.PortSelector = v.PortSelector
		rule.Ingress = v.Ingress

		if _, err := rule.Validate(nil); err != nil {
			return nil, err
		}

		result = append(result, rule)
	}

	return result, nil
}

// marshalYaml encodes `in` into `yaml` bytes with 2 indentation.
// It also returns an error, if any.
func marshalYaml(in any) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(in); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// CombineYamlBytes prepends and returns `---\n` before `input`
func CombineYamlBytes(input [][]byte) []byte {
	delimiter := []byte("---\n")
	var result []byte
	for k := range input {
		// https://github.com/budimanjojo/talhelper/issues/497
		if !bytes.HasPrefix(input[k], delimiter) {
			result = append(result, delimiter...)
		}
		result = append(result, input[k]...)
	}
	return result
}
