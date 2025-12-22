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

	addCommonLinkConfig(&bondConfig.CommonLinkConfig, device)

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

func GenerateAddressConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if len(device.DeviceAddresses) == 0 || hasSpecialConfig(device) {
			continue
		}

		addressConfig := GenerateAddressConfig(device)
		if addressConfig == nil {
			continue
		}

		addressBytes, err := marshalYaml(addressConfig)
		if err != nil {
			return nil, err
		}

		result = append(result, addressBytes)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return CombineYamlBytes(result), nil
}

func GenerateAddressConfig(device *v1alpha1.Device) *network.LinkConfigV1Alpha1 {
	if device == nil || len(device.DeviceAddresses) == 0 {
		return nil
	}

	linkConfig := network.NewLinkConfigV1Alpha1(device.DeviceInterface)

	for _, address := range device.DeviceAddresses {
		prefix, err := netip.ParsePrefix(address)
		if err != nil {
			ip, ipErr := netip.ParseAddr(address)
			if ipErr != nil {
				continue
			}
			bits := 32
			if ip.Is6() {
				bits = 128
			}
			prefix = netip.PrefixFrom(ip, bits)
		}

		linkConfig.LinkAddresses = append(linkConfig.LinkAddresses, network.AddressConfig{
			AddressAddress: prefix,
		})
	}

	return linkConfig
}

func GenerateRouteConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if len(device.DeviceRoutes) == 0 || hasSpecialConfig(device) {
			continue
		}

		routeConfig := GenerateRouteConfig(device)
		if routeConfig == nil {
			continue
		}

		routeBytes, err := marshalYaml(routeConfig)
		if err != nil {
			return nil, err
		}

		result = append(result, routeBytes)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return CombineYamlBytes(result), nil
}

func GenerateRouteConfig(device *v1alpha1.Device) *network.LinkConfigV1Alpha1 {
	if device == nil || len(device.DeviceRoutes) == 0 {
		return nil
	}

	linkConfig := network.NewLinkConfigV1Alpha1(device.DeviceInterface)

	for _, route := range device.DeviceRoutes {
		routeConfig := network.RouteConfig{}

		if route.Network() != "" {
			prefix, err := netip.ParsePrefix(route.Network())
			if err == nil {
				routeConfig.RouteDestination = network.Prefix{Prefix: prefix}
			}
		}

		if route.Gateway() != "" {
			gateway, err := netip.ParseAddr(route.Gateway())
			if err == nil {
				routeConfig.RouteGateway = network.Addr{Addr: gateway}
			}
		}

		if route.Source() != "" {
			source, err := netip.ParseAddr(route.Source())
			if err == nil {
				routeConfig.RouteSource = network.Addr{Addr: source}
			}
		}

		if route.Metric() > 0 {
			routeConfig.RouteMetric = route.Metric()
		}

		if route.MTU() > 0 {
			routeConfig.RouteMTU = route.MTU()
		}

		linkConfig.LinkRoutes = append(linkConfig.LinkRoutes, routeConfig)
	}

	return linkConfig
}

func GenerateLinkConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if hasSpecialConfig(device) {
			continue
		}
		linkConfig := GenerateLinkConfig(device)
		if linkConfig == nil {
			continue
		}

		linkBytes, err := marshalYaml(linkConfig)
		if err != nil {
			return nil, err
		}

		result = append(result, linkBytes)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return CombineYamlBytes(result), nil
}

func GenerateLinkConfig(device *v1alpha1.Device) *network.LinkConfigV1Alpha1 {
	if device == nil {
		return nil
	}

	if device.DeviceMTU == 0 {
		return nil
	}

	linkConfig := network.NewLinkConfigV1Alpha1(device.DeviceInterface)

	if device.DeviceMTU > 0 {
		linkConfig.LinkMTU = uint32(device.DeviceMTU)
	}

	return linkConfig
}

func GenerateVLANConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if len(device.DeviceVlans) == 0 {
			continue
		}

		for _, vlan := range device.DeviceVlans {
			vlanConfig := GenerateVLANConfig(device, vlan)
			if vlanConfig == nil {
				continue
			}

			vlanBytes, err := marshalYaml(vlanConfig)
			if err != nil {
				return nil, err
			}

			result = append(result, vlanBytes)
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return CombineYamlBytes(result), nil
}

func GenerateVLANConfig(device *v1alpha1.Device, vlan *v1alpha1.Vlan) *network.VLANConfigV1Alpha1 {
	if device == nil || vlan == nil {
		return nil
	}

	vlanInterface := device.DeviceInterface
	if vlan.VlanID > 0 {
		vlanConfig := network.NewVLANConfigV1Alpha1(vlanInterface)
		vlanConfig.VLANIDConfig = vlan.VlanID
		vlanConfig.ParentLinkConfig = device.DeviceInterface

		if len(vlan.VlanAddresses) > 0 {
			for _, addr := range vlan.VlanAddresses {
				prefix, err := netip.ParsePrefix(addr)
				if err != nil {
					ip, ipErr := netip.ParseAddr(addr)
					if ipErr != nil {
						continue
					}
					bits := 32
					if ip.Is6() {
						bits = 128
					}
					prefix = netip.PrefixFrom(ip, bits)
				}
				vlanConfig.LinkAddresses = append(vlanConfig.LinkAddresses, network.AddressConfig{
					AddressAddress: prefix,
				})
			}
		}

		if len(vlan.VlanRoutes) > 0 {
			for _, route := range vlan.VlanRoutes {
				routeSpec := network.RouteConfig{}

				if route.Network() != "" {
					prefix, err := netip.ParsePrefix(route.Network())
					if err == nil {
						routeSpec.RouteDestination = network.Prefix{Prefix: prefix}
					}
				}

				if route.Gateway() != "" {
					gateway, err := netip.ParseAddr(route.Gateway())
					if err == nil {
						routeSpec.RouteGateway = network.Addr{Addr: gateway}
					}
				}

				if route.Source() != "" {
					source, err := netip.ParseAddr(route.Source())
					if err == nil {
						routeSpec.RouteSource = network.Addr{Addr: source}
					}
				}

				if route.Metric() > 0 {
					routeSpec.RouteMetric = route.Metric()
				}

				if route.MTU() > 0 {
					routeSpec.RouteMTU = route.MTU()
				}

				vlanConfig.LinkRoutes = append(vlanConfig.LinkRoutes, routeSpec)
			}
		}

		if vlan.VlanMTU > 0 {
			vlanConfig.LinkMTU = uint32(vlan.VlanMTU)
		}

		return vlanConfig
	}

	return nil
}

func GenerateWireguardConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if device.DeviceWireguardConfig == nil {
			continue
		}

		wgConfig := GenerateWireguardConfig(device)
		if wgConfig == nil {
			continue
		}

		wgBytes, err := marshalYaml(wgConfig)
		if err != nil {
			return nil, err
		}

		result = append(result, wgBytes)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return CombineYamlBytes(result), nil
}

func GenerateWireguardConfig(device *v1alpha1.Device) *network.WireguardConfigV1Alpha1 {
	if device == nil || device.DeviceWireguardConfig == nil {
		return nil
	}

	wgConfig := network.NewWireguardConfigV1Alpha1(device.DeviceInterface)

	if device.DeviceWireguardConfig.WireguardPrivateKey != "" {
		wgConfig.WireguardPrivateKey = device.DeviceWireguardConfig.WireguardPrivateKey
	}

	if device.DeviceWireguardConfig.WireguardListenPort > 0 {
		wgConfig.WireguardListenPort = device.DeviceWireguardConfig.WireguardListenPort
	}

	if device.DeviceWireguardConfig.WireguardFirewallMark > 0 {
		wgConfig.WireguardFirewallMark = device.DeviceWireguardConfig.WireguardFirewallMark
	}

	if len(device.DeviceWireguardConfig.WireguardPeers) > 0 {
		for _, peer := range device.DeviceWireguardConfig.WireguardPeers {
			wgPeer := network.WireguardPeer{}

			if peer.WireguardPublicKey != "" {
				wgPeer.WireguardPublicKey = peer.WireguardPublicKey
			}

			if peer.WireguardEndpoint != "" {
				addrPort, err := netip.ParseAddrPort(peer.WireguardEndpoint)
				if err == nil {
					wgPeer.WireguardEndpoint = network.AddrPort{AddrPort: addrPort}
				}
			}

			if peer.WireguardPersistentKeepaliveInterval > 0 {
				wgPeer.WireguardPersistentKeepaliveInterval = peer.WireguardPersistentKeepaliveInterval
			}

			if len(peer.WireguardAllowedIPs) > 0 {
				for _, allowedIP := range peer.WireguardAllowedIPs {
					prefix, err := netip.ParsePrefix(allowedIP)
					if err == nil {
						wgPeer.WireguardAllowedIPs = append(wgPeer.WireguardAllowedIPs, network.Prefix{Prefix: prefix})
					}
				}
			}

			wgConfig.WireguardPeers = append(wgConfig.WireguardPeers, wgPeer)
		}
	}

	addCommonLinkConfig(&wgConfig.CommonLinkConfig, device)

	return wgConfig
}

func GenerateBridgeConfigBytes(devices []*v1alpha1.Device) ([]byte, error) {
	var result [][]byte

	for _, device := range devices {
		if device.DeviceBridge == nil {
			continue
		}

		bridgeConfig := GenerateBridgeConfig(device)
		if bridgeConfig == nil {
			continue
		}

		bridgeBytes, err := marshalYaml(bridgeConfig)
		if err != nil {
			return nil, err
		}

		result = append(result, bridgeBytes)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return CombineYamlBytes(result), nil
}

func GenerateBridgeConfig(device *v1alpha1.Device) *network.BridgeConfigV1Alpha1 {
	if device == nil || device.DeviceBridge == nil {
		return nil
	}

	bridgeConfig := network.NewBridgeConfigV1Alpha1(device.DeviceInterface)

	if len(device.DeviceBridge.BridgedInterfaces) > 0 {
		bridgeConfig.BridgeLinks = device.DeviceBridge.BridgedInterfaces
	}

	if device.DeviceBridge.BridgeSTP != nil && device.DeviceBridge.BridgeSTP.STPEnabled != nil {
		bridgeConfig.BridgeSTP = network.BridgeSTPConfig{
			BridgeSTPEnabled: device.DeviceBridge.BridgeSTP.STPEnabled,
		}
	}

	addCommonLinkConfig(&bridgeConfig.CommonLinkConfig, device)

	return bridgeConfig
}

func hasSpecialConfig(device *v1alpha1.Device) bool {
	if device == nil {
		return false
	}
	return device.DeviceBond != nil || len(device.DeviceVlans) > 0 ||
		device.DeviceWireguardConfig != nil || device.DeviceBridge != nil
}

func addCommonLinkConfig(linkConfig *network.CommonLinkConfig, device *v1alpha1.Device) {
	if device == nil || linkConfig == nil {
		return
	}

	for _, address := range device.DeviceAddresses {
		prefix, err := netip.ParsePrefix(address)
		if err != nil {
			ip, ipErr := netip.ParseAddr(address)
			if ipErr != nil {
				continue
			}
			bits := 32
			if ip.Is6() {
				bits = 128
			}
			prefix = netip.PrefixFrom(ip, bits)
		}

		linkConfig.LinkAddresses = append(linkConfig.LinkAddresses, network.AddressConfig{
			AddressAddress: prefix,
		})
	}

	for _, route := range device.DeviceRoutes {
		routeConfig := network.RouteConfig{}

		if route.Network() != "" {
			prefix, err := netip.ParsePrefix(route.Network())
			if err == nil {
				routeConfig.RouteDestination = network.Prefix{Prefix: prefix}
			}
		}

		if route.Gateway() != "" {
			gateway, err := netip.ParseAddr(route.Gateway())
			if err == nil {
				routeConfig.RouteGateway = network.Addr{Addr: gateway}
			}
		}

		if route.Source() != "" {
			source, err := netip.ParseAddr(route.Source())
			if err == nil {
				routeConfig.RouteSource = network.Addr{Addr: source}
			}
		}

		if route.Metric() > 0 {
			routeConfig.RouteMetric = route.Metric()
		}

		if route.MTU() > 0 {
			routeConfig.RouteMTU = route.MTU()
		}

		linkConfig.LinkRoutes = append(linkConfig.LinkRoutes, routeConfig)
	}

	if device.DeviceMTU > 0 {
		linkConfig.LinkMTU = uint32(device.DeviceMTU)
	}
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
