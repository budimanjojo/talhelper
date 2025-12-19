package talos

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/go-pointer"
	"github.com/siderolabs/talos/pkg/machinery/config/types/network"
	"github.com/siderolabs/talos/pkg/machinery/nethelpers"
	"gopkg.in/yaml.v3"
)

func TestGenerateNetworkHostname(t *testing.T) {
	result1 := GenerateNetworkHostnameConfig("shouldbeignored", true)
	result2 := GenerateNetworkHostnameConfig("hostname", false)

	compare(result1.ConfigAuto, pointer.To(nethelpers.AutoHostnameKindStable), t)
	compare(result1.ConfigHostname, "", t)
	compare(result2.ConfigAuto, pointer.To(nethelpers.AutoHostnameKindOff), t)
	compare(result2.ConfigHostname, "hostname", t)
}

func TestGenerateNodeDefaultActionConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    ingressFirewall:
      defaultAction: accept
  - hostname: node2
    ingressFirewall:
      defaultAction: block`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	node1Result := GenerateNodeDefaultActionConfig(m.Nodes[0].IngressFirewall)
	node2Result := GenerateNodeDefaultActionConfig(m.Nodes[1].IngressFirewall)

	compare(node1Result.Ingress.String(), "accept", t)
	compare(node2Result.Ingress.String(), "block", t)
}

func TestGenerateNodeRuleConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    ingressFirewall:
      defaultAction: accept
      rules:
        - name: kubelet-ingress
          portSelector:
            ports:
              - 10250
            protocol: tcp
          ingress:
            - subnet: 172.20.0.0/24
              except: 172.20.0.1/32
        - name: etcd-ingress
          portSelector:
            ports:
              - 2379-2380
            protocol: tcp
          ingress:
            - subnet: 10.10.10.1/32
            - subnet: 10.10.10.2/32
            - subnet: 10.10.10.3/32`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	expectedRule1Name := "kubelet-ingress"
	expectedRule1PortSelector := network.RulePortSelector{
		Ports:    network.PortRanges{network.PortRange{Lo: 10250, Hi: 10250}},
		Protocol: nethelpers.ProtocolTCP,
	}
	expectedRule1Ingress := network.IngressConfig{
		{
			Subnet: netip.MustParsePrefix("172.20.0.0/24"),
			Except: network.Prefix{Prefix: netip.MustParsePrefix("172.20.0.1/32")},
		},
	}
	expectedRule2Name := "etcd-ingress"
	expectedRule2PortSelector := network.RulePortSelector{
		Ports:    network.PortRanges{network.PortRange{Lo: 2379, Hi: 2380}},
		Protocol: nethelpers.ProtocolTCP,
	}
	expectedRule2Ingress := network.IngressConfig{
		{Subnet: netip.MustParsePrefix("10.10.10.1/32")},
		{Subnet: netip.MustParsePrefix("10.10.10.2/32")},
		{Subnet: netip.MustParsePrefix("10.10.10.3/32")},
	}

	result, err := GenerateNodeRuleConfig(m.Nodes[0].IngressFirewall)
	if err != nil {
		t.Fatal(err)
	}

	compare(result[0].Name(), expectedRule1Name, t)
	compare(result[0].PortSelector, expectedRule1PortSelector, t)
	compare(result[0].Ingress, expectedRule1Ingress, t)
	compare(result[1].Name(), expectedRule2Name, t)
	compare(result[1].PortSelector, expectedRule2PortSelector, t)
	compare(result[1].Ingress, expectedRule2Ingress, t)
}

func TestGenerateBondConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: bond0
        bond:
          interfaces:
            - enp1s2
            - enp1s3
          mode: 802.3ad
          lacpRate: fast
          xmitHashPolicy: layer3+4
          miimon: 100
          updelay: 200
          downdelay: 200
          arpIPTarget:
            - 10.15.0.1
        mtu: 9000`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if len(m.Nodes) == 0 || len(m.Nodes[0].NetworkInterfaces) == 0 {
		t.Fatal("failed to parse test data")
	}

	result := GenerateBondConfig(m.Nodes[0].NetworkInterfaces[0])
	if result == nil {
		t.Fatal("expected bond config, got nil")
	}

	if result.Name() != "bond0" {
		t.Errorf("expected name=bond0, got %s", result.Name())
	}

	if len(result.BondLinks) != 2 {
		t.Errorf("expected 2 bond links, got %d", len(result.BondLinks))
	}
	if result.BondLinks[0] != "enp1s2" {
		t.Errorf("expected link enp1s2, got %s", result.BondLinks[0])
	}
	if result.BondLinks[1] != "enp1s3" {
		t.Errorf("expected link enp1s3, got %s", result.BondLinks[1])
	}

	if result.BondMode == nil {
		t.Error("expected bond mode to be set")
	} else if result.BondMode.String() != "802.3ad" {
		t.Errorf("expected mode=802.3ad, got %s", result.BondMode.String())
	}

	if result.BondLACPRate == nil {
		t.Error("expected LACP rate to be set")
	} else if result.BondLACPRate.String() != "fast" {
		t.Errorf("expected LACP rate=fast, got %s", result.BondLACPRate.String())
	}

	if result.BondXmitHashPolicy == nil {
		t.Error("expected xmit hash policy to be set")
	} else if result.BondXmitHashPolicy.String() != "layer3+4" {
		t.Errorf("expected xmit hash policy=layer3+4, got %s", result.BondXmitHashPolicy.String())
	}

	if result.BondMIIMon == nil || *result.BondMIIMon != 100 {
		t.Errorf("expected miimon=100, got %v", result.BondMIIMon)
	}
	if result.BondUpDelay == nil || *result.BondUpDelay != 200 {
		t.Errorf("expected updelay=200, got %v", result.BondUpDelay)
	}
	if result.BondDownDelay == nil || *result.BondDownDelay != 200 {
		t.Errorf("expected downdelay=200, got %v", result.BondDownDelay)
	}

	if len(result.BondARPIPTargets) != 1 {
		t.Errorf("expected 1 ARP target, got %d", len(result.BondARPIPTargets))
	} else if result.BondARPIPTargets[0].String() != "10.15.0.1" {
		t.Errorf("expected ARP target=10.15.0.1, got %s", result.BondARPIPTargets[0].String())
	}
}

func TestGenerateDHCP4Config(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        dhcp: true
        dhcpOptions:
          routeMetric: 1024`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if len(m.Nodes) == 0 || len(m.Nodes[0].NetworkInterfaces) == 0 {
		t.Fatal("failed to parse test data")
	}

	result := GenerateDHCP4Config(m.Nodes[0].NetworkInterfaces[0])
	if result == nil {
		t.Fatal("expected DHCP4 config, got nil")
	}

	if result.Name() != "eth0" {
		t.Errorf("expected name=eth0, got %s", result.Name())
	}

	if result.ConfigRouteMetric != 1024 {
		t.Errorf("expected route metric=1024, got %d", result.ConfigRouteMetric)
	}
}

func TestGenerateDHCP6Config(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        dhcpOptions:
          ipv6: true
          routeMetric: 2048`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if len(m.Nodes) == 0 || len(m.Nodes[0].NetworkInterfaces) == 0 {
		t.Fatal("failed to parse test data")
	}

	result := GenerateDHCP6Config(m.Nodes[0].NetworkInterfaces[0])
	if result == nil {
		t.Fatal("expected DHCP6 config, got nil")
	}

	if result.Name() != "eth0" {
		t.Errorf("expected name=eth0, got %s", result.Name())
	}

	if result.ConfigRouteMetric != 2048 {
		t.Errorf("expected route metric=2048, got %d", result.ConfigRouteMetric)
	}
}

func TestGenerateBondConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: bond0
        bond:
          interfaces:
            - enp1s2
          mode: balance-rr
      - interface: eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	bondBytes, err := GenerateBondConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if bondBytes == nil {
		t.Fatal("expected bond config bytes, got nil")
	}

	bondStr := string(bondBytes)
	if !bytes.Contains(bondBytes, []byte("kind: BondConfig")) {
		t.Error("expected output to contain 'kind: BondConfig'")
	}
	if !bytes.Contains(bondBytes, []byte("name: bond0")) {
		t.Error("expected output to contain 'name: bond0'")
	}
	if !bytes.Contains(bondBytes, []byte("---")) {
		t.Error("expected output to contain YAML document delimiter")
	}
	t.Logf("Bond config output:\n%s", bondStr)
}

func TestGenerateDHCP4ConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        dhcp: true
      - interface: eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	dhcpBytes, err := GenerateDHCP4ConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if dhcpBytes == nil {
		t.Fatal("expected DHCP4 config bytes, got nil")
	}

	dhcpStr := string(dhcpBytes)
	if !bytes.Contains(dhcpBytes, []byte("kind: DHCPv4Config")) {
		t.Error("expected output to contain 'kind: DHCPv4Config'")
	}
	if !bytes.Contains(dhcpBytes, []byte("name: eth0")) {
		t.Error("expected output to contain 'name: eth0'")
	}
	t.Logf("DHCP4 config output:\n%s", dhcpStr)
}

func TestGenerateVIPConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        vip:
          ip: 192.168.1.100`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if len(m.Nodes) == 0 || len(m.Nodes[0].NetworkInterfaces) == 0 {
		t.Fatal("failed to parse test data")
	}

	result := GenerateVIPConfig(m.Nodes[0].NetworkInterfaces[0])
	if result == nil {
		t.Fatal("expected VIP config, got nil")
	}

	if result.Name() != "192.168.1.100" {
		t.Errorf("expected VIP=192.168.1.100, got %s", result.Name())
	}

	if result.LinkName != "eth0" {
		t.Errorf("expected link=eth0, got %s", result.LinkName)
	}
}

func TestGenerateVIPConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        vip:
          ip: 192.168.1.100
      - interface: eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	vipBytes, err := GenerateVIPConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if vipBytes == nil {
		t.Fatal("expected VIP config bytes, got nil")
	}

	vipStr := string(vipBytes)
	if !bytes.Contains(vipBytes, []byte("kind: Layer2VIPConfig")) {
		t.Error("expected output to contain 'kind: Layer2VIPConfig'")
	}
	if !bytes.Contains(vipBytes, []byte("name: 192.168.1.100")) {
		t.Error("expected output to contain 'name: 192.168.1.100'")
	}
	if !bytes.Contains(vipBytes, []byte("link: eth0")) {
		t.Error("expected output to contain 'link: eth0'")
	}
	t.Logf("VIP config output:\n%s", vipStr)
}
