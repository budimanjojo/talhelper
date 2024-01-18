package talos

import (
	"bytes"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/types/network"
	"gopkg.in/yaml.v3"
)

func GenerateNodeNetworkConfigBytes(node *config.Node) ([]byte, error) {
	var result [][]byte

	defaultAction := GenerateNodeDefaultActionConfig(node)
	defaultActionBytes, err := marshalYaml(defaultAction)
	if err != nil {
		return nil, err
	}

	result = append(result, defaultActionBytes)

	rules, err := GenerateNodeRuleConfig(node)
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

func GenerateNodeDefaultActionConfig(node *config.Node) *network.DefaultActionConfigV1Alpha1 {
	result := network.NewDefaultActionConfigV1Alpha1()
	result.Ingress = node.IngressFirewall.DefaultAction

	return result
}

func GenerateNodeRuleConfig(node *config.Node) ([]*network.RuleConfigV1Alpha1, error) {
	var result []*network.RuleConfigV1Alpha1

	for _, v := range node.IngressFirewall.NetworkRules {
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
func marshalYaml(in interface{}) ([]byte, error) {
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
		result = append(result, delimiter...)
		result = append(result, input[k]...)
	}
	return result
}
