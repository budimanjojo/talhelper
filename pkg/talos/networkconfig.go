package talos

import (
	"bytes"
	"net/netip"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/go-pointer"
	"github.com/siderolabs/talos/pkg/machinery/config/types/network"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/siderolabs/talos/pkg/machinery/nethelpers"
	"gopkg.in/yaml.v3"
)

func GenerateResolverConfigBytes(nameservers []string, disableSearchDomain bool) ([]byte, error) {
	var result [][]byte

	resolver := GenerateResolverConfig(nameservers, disableSearchDomain)
	resolverBytes, err := marshalYaml(resolver)
	if err != nil {
		return nil, err
	}

	return CombineYamlBytes(append(result, resolverBytes)), nil
}

func GenerateResolverConfig(nameservers []string, disableSearchDomain bool) *network.ResolverConfigV1Alpha1 {
	result := network.NewResolverConfigV1Alpha1()

	ns := []network.NameserverConfig{}
	for _, n := range nameservers {
		ns = append(ns, network.NameserverConfig{Address: network.Addr{Addr: netip.MustParseAddr(n)}})
	}

	result.ResolverNameservers = ns
	result.ResolverSearchDomains.SearchDisableDefault = &disableSearchDomain

	return result
}

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

func GenerateBondConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if device.DeviceBond == nil {
			continue
		}

		bondConfig := GenerateBondConfig(device)
		if bondConfig == nil {
			continue
		}

		bondBytes, err := marshalYaml(bondConfig)
		if err != nil {
			return nil, err
		}

		result = append(result, bondBytes)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return CombineYamlBytes(result), nil
}

func GenerateBondConfig(device *v1alpha1.Device) *network.BondConfigV1Alpha1 {
	if device == nil || device.DeviceBond == nil {
		return nil
	}

	bondConfig := network.NewBondConfigV1Alpha1(device.DeviceInterface)

	if len(device.DeviceBond.BondInterfaces) > 0 {
		bondConfig.BondLinks = device.DeviceBond.BondInterfaces
	}

	if device.DeviceBond.BondMode != "" {
		if mode, err := nethelpers.BondModeByName(device.DeviceBond.BondMode); err == nil {
			bondConfig.BondMode = &mode
		}
	}

	if device.DeviceBond.BondHashPolicy != "" {
		if policy, err := nethelpers.BondXmitHashPolicyByName(device.DeviceBond.BondHashPolicy); err == nil {
			bondConfig.BondXmitHashPolicy = &policy
		}
	}

	if device.DeviceBond.BondLACPRate != "" {
		if rate, err := nethelpers.LACPRateByName(device.DeviceBond.BondLACPRate); err == nil {
			bondConfig.BondLACPRate = &rate
		}
	}

	if device.DeviceBond.BondARPValidate != "" {
		if validate, err := nethelpers.ARPValidateByName(device.DeviceBond.BondARPValidate); err == nil {
			bondConfig.BondARPValidate = &validate
		}
	}

	if device.DeviceBond.BondARPAllTargets != "" {
		if targets, err := nethelpers.ARPAllTargetsByName(device.DeviceBond.BondARPAllTargets); err == nil {
			bondConfig.BondARPAllTargets = &targets
		}
	}

	if device.DeviceBond.BondPrimaryReselect != "" {
		if reselect, err := nethelpers.PrimaryReselectByName(device.DeviceBond.BondPrimaryReselect); err == nil {
			bondConfig.BondPrimaryReselect = &reselect
		}
	}

	if device.DeviceBond.BondFailOverMac != "" {
		if failOver, err := nethelpers.FailOverMACByName(device.DeviceBond.BondFailOverMac); err == nil {
			bondConfig.BondFailOverMAC = &failOver
		}
	}

	if device.DeviceBond.BondADSelect != "" {
		if adSelect, err := nethelpers.ADSelectByName(device.DeviceBond.BondADSelect); err == nil {
			bondConfig.BondADSelect = &adSelect
		}
	}

	if device.DeviceBond.BondMIIMon > 0 {
		bondConfig.BondMIIMon = pointer.To(device.DeviceBond.BondMIIMon)
	}
	if device.DeviceBond.BondUpDelay > 0 {
		bondConfig.BondUpDelay = pointer.To(device.DeviceBond.BondUpDelay)
	}
	if device.DeviceBond.BondDownDelay > 0 {
		bondConfig.BondDownDelay = pointer.To(device.DeviceBond.BondDownDelay)
	}
	if device.DeviceBond.BondARPInterval > 0 {
		bondConfig.BondARPInterval = pointer.To(device.DeviceBond.BondARPInterval)
	}
	if device.DeviceBond.BondResendIGMP > 0 {
		bondConfig.BondResendIGMP = pointer.To(device.DeviceBond.BondResendIGMP)
	}
	if device.DeviceBond.BondMinLinks > 0 {
		bondConfig.BondMinLinks = pointer.To(device.DeviceBond.BondMinLinks)
	}
	if device.DeviceBond.BondLPInterval > 0 {
		bondConfig.BondLPInterval = pointer.To(device.DeviceBond.BondLPInterval)
	}
	if device.DeviceBond.BondPacketsPerSlave > 0 {
		bondConfig.BondPacketsPerSlave = pointer.To(device.DeviceBond.BondPacketsPerSlave)
	}
	if device.DeviceBond.BondNumPeerNotif > 0 {
		bondConfig.BondNumPeerNotif = pointer.To(device.DeviceBond.BondNumPeerNotif)
	}
	if device.DeviceBond.BondTLBDynamicLB > 0 {
		bondConfig.BondTLBDynamicLB = pointer.To(device.DeviceBond.BondTLBDynamicLB)
	}
	if device.DeviceBond.BondAllSlavesActive > 0 {
		bondConfig.BondAllSlavesActive = pointer.To(device.DeviceBond.BondAllSlavesActive)
	}
	if device.DeviceBond.BondUseCarrier != nil {
		bondConfig.BondUseCarrier = device.DeviceBond.BondUseCarrier
	}
	if device.DeviceBond.BondADActorSysPrio > 0 {
		bondConfig.BondADActorSysPrio = pointer.To(device.DeviceBond.BondADActorSysPrio)
	}
	if device.DeviceBond.BondADUserPortKey > 0 {
		bondConfig.BondADUserPortKey = pointer.To(device.DeviceBond.BondADUserPortKey)
	}
	if device.DeviceBond.BondPeerNotifyDelay > 0 {
		bondConfig.BondPeerNotifyDelay = pointer.To(device.DeviceBond.BondPeerNotifyDelay)
	}

	if len(device.DeviceBond.BondARPIPTarget) > 0 {
		bondConfig.BondARPIPTargets = make([]netip.Addr, len(device.DeviceBond.BondARPIPTarget))
		for i, ip := range device.DeviceBond.BondARPIPTarget {
			bondConfig.BondARPIPTargets[i] = netip.MustParseAddr(ip)
		}
	}

	return bondConfig
}

func GenerateDHCP4ConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if device.DeviceDHCP == nil || !*device.DeviceDHCP {
			continue
		}

		dhcpConfig := GenerateDHCP4Config(device)
		if dhcpConfig == nil {
			continue
		}

		dhcpBytes, err := marshalYaml(dhcpConfig)
		if err != nil {
			return nil, err
		}

		result = append(result, dhcpBytes)
	}

	if len(result) == 0 {
		return nil, nil
	}
	return CombineYamlBytes(result), nil
}

func GenerateDHCP4Config(device *v1alpha1.Device) *network.DHCPv4ConfigV1Alpha1 {
	if device == nil || device.DeviceDHCP == nil || !*device.DeviceDHCP {
		return nil
	}

	dhcpConfig := network.NewDHCPv4ConfigV1Alpha1(device.DeviceInterface)

	if device.DeviceDHCPOptions != nil {
		if device.DeviceDHCPOptions.DHCPRouteMetric != 0 {
			dhcpConfig.ConfigRouteMetric = device.DeviceDHCPOptions.DHCPRouteMetric
		}
	}

	return dhcpConfig
}

func GenerateDHCP6ConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if device.DeviceDHCPOptions == nil || device.DeviceDHCPOptions.DHCPIPv6 == nil || !*device.DeviceDHCPOptions.DHCPIPv6 {
			continue
		}

		dhcp6Config := GenerateDHCP6Config(device)
		if dhcp6Config == nil {
			continue
		}

		dhcp6Bytes, err := marshalYaml(dhcp6Config)
		if err != nil {
			return nil, err
		}

		result = append(result, dhcp6Bytes)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return CombineYamlBytes(result), nil
}

func GenerateDHCP6Config(device *v1alpha1.Device) *network.DHCPv6ConfigV1Alpha1 {
	if device == nil || device.DeviceDHCPOptions == nil || device.DeviceDHCPOptions.DHCPIPv6 == nil || !*device.DeviceDHCPOptions.DHCPIPv6 {
		return nil
	}

	dhcp6Config := network.NewDHCPv6ConfigV1Alpha1(device.DeviceInterface)

	if device.DeviceDHCPOptions.DHCPRouteMetric != 0 {
		dhcp6Config.ConfigRouteMetric = device.DeviceDHCPOptions.DHCPRouteMetric
	}

	return dhcp6Config
}

func GenerateVIPConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if device.DeviceVIPConfig == nil || device.DeviceVIPConfig.SharedIP == "" {
			continue
		}

		vipConfig := GenerateVIPConfig(device)
		if vipConfig == nil {
			continue
		}

		vipBytes, err := marshalYaml(vipConfig)
		if err != nil {
			return nil, err
		}

		result = append(result, vipBytes)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return CombineYamlBytes(result), nil
}

func GenerateVIPConfig(device *v1alpha1.Device) *network.Layer2VIPConfigV1Alpha1 {
	if device == nil || device.DeviceVIPConfig == nil || device.DeviceVIPConfig.SharedIP == "" {
		return nil
	}

	vipConfig := network.NewLayer2VIPConfigV1Alpha1(device.DeviceVIPConfig.SharedIP)
	vipConfig.LinkName = device.DeviceInterface

	return vipConfig
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
