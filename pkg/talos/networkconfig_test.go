package talos

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/go-pointer"
	netenums "github.com/siderolabs/talos/pkg/machinery/api/resource/definitions/enums"
	networktypes "github.com/siderolabs/talos/pkg/machinery/api/resource/definitions/network"
	"github.com/siderolabs/talos/pkg/machinery/cel"
	"github.com/siderolabs/talos/pkg/machinery/cel/celenv"
	"github.com/siderolabs/talos/pkg/machinery/config/types/network"
	v1alpha1 "github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
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

func TestGenerateAddressConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        addresses:
          - 192.168.1.100/24
          - 10.0.0.1
          - 2001:db8::1/64`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	result := GenerateLinkConfig(m.Nodes[0].NetworkInterfaces[0])
	if result == nil {
		t.Fatal("expected link config, got nil")
	}

	if result.MetaName != "eth0" {
		t.Errorf("expected name=eth0, got %s", result.MetaName)
	}

	if len(result.LinkAddresses) != 3 {
		t.Fatalf("expected 3 addresses, got %d", len(result.LinkAddresses))
	}

	expectedAddr1 := netip.MustParsePrefix("192.168.1.100/24")
	if result.LinkAddresses[0].AddressAddress != expectedAddr1 {
		t.Errorf("expected address %s, got %s", expectedAddr1, result.LinkAddresses[0].AddressAddress)
	}

	expectedAddr2 := netip.MustParsePrefix("10.0.0.1/32")
	if result.LinkAddresses[1].AddressAddress != expectedAddr2 {
		t.Errorf("expected address %s, got %s", expectedAddr2, result.LinkAddresses[1].AddressAddress)
	}

	expectedAddr3 := netip.MustParsePrefix("2001:db8::1/64")
	if result.LinkAddresses[2].AddressAddress != expectedAddr3 {
		t.Errorf("expected address %s, got %s", expectedAddr3, result.LinkAddresses[2].AddressAddress)
	}
}

func TestGenerateAddressConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        addresses:
          - 192.168.1.100/24
      - interface: eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	linkBytes, err := GenerateLinkConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if linkBytes == nil {
		t.Fatal("expected link config bytes, got nil")
	}

	linkStr := string(linkBytes)
	if !bytes.Contains(linkBytes, []byte("kind: LinkConfig")) {
		t.Error("expected output to contain 'kind: LinkConfig'")
	}
	if !bytes.Contains(linkBytes, []byte("name: eth0")) {
		t.Error("expected output to contain 'name: eth0'")
	}
	if !bytes.Contains(linkBytes, []byte("address: 192.168.1.100/24")) {
		t.Error("expected output to contain 'address: 192.168.1.100/24'")
	}
	t.Logf("Link config output:\n%s", linkStr)
}

func TestGenerateRouteConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        routes:
          - network: 10.0.0.0/8
            gateway: 192.168.1.1
            metric: 100
          - network: 0.0.0.0/0
            gateway: 192.168.1.254`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	result := GenerateLinkConfig(m.Nodes[0].NetworkInterfaces[0])
	if result == nil {
		t.Fatal("expected link config, got nil")
	}

	if result.MetaName != "eth0" {
		t.Errorf("expected name=eth0, got %s", result.MetaName)
	}

	if len(result.LinkRoutes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(result.LinkRoutes))
	}

	expectedDest1 := netip.MustParsePrefix("10.0.0.0/8")
	if result.LinkRoutes[0].RouteDestination.Prefix != expectedDest1 {
		t.Errorf("expected destination %s, got %s", expectedDest1, result.LinkRoutes[0].RouteDestination.Prefix)
	}
	expectedGw1 := netip.MustParseAddr("192.168.1.1")
	if result.LinkRoutes[0].RouteGateway.Addr != expectedGw1 {
		t.Errorf("expected gateway %s, got %s", expectedGw1, result.LinkRoutes[0].RouteGateway.Addr)
	}
	if result.LinkRoutes[0].RouteMetric != 100 {
		t.Errorf("expected metric 100, got %d", result.LinkRoutes[0].RouteMetric)
	}

	// For default routes (0.0.0.0/0), the destination field should be omitted (zero value)
	// so Talos can infer it from the gateway's address family
	if result.LinkRoutes[1].RouteDestination.IsValid() {
		t.Errorf("expected no destination for default route, got %s", result.LinkRoutes[1].RouteDestination.Prefix)
	}
	expectedGw2 := netip.MustParseAddr("192.168.1.254")
	if result.LinkRoutes[1].RouteGateway.Addr != expectedGw2 {
		t.Errorf("expected gateway %s, got %s", expectedGw2, result.LinkRoutes[1].RouteGateway.Addr)
	}
}

func TestGenerateRouteConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        routes:
          - network: 10.0.0.0/8
            gateway: 192.168.1.1
      - interface: eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	linkBytes, err := GenerateLinkConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if linkBytes == nil {
		t.Fatal("expected link config bytes, got nil")
	}

	linkStr := string(linkBytes)
	if !bytes.Contains(linkBytes, []byte("kind: LinkConfig")) {
		t.Error("expected output to contain 'kind: LinkConfig'")
	}
	if !bytes.Contains(linkBytes, []byte("name: eth0")) {
		t.Error("expected output to contain 'name: eth0'")
	}
	if !bytes.Contains(linkBytes, []byte("destination: 10.0.0.0/8")) {
		t.Error("expected output to contain 'destination: 10.0.0.0/8'")
	}
	if !bytes.Contains(linkBytes, []byte("gateway: 192.168.1.1")) {
		t.Error("expected output to contain 'gateway: 192.168.1.1'")
	}
	t.Logf("Link config output:\n%s", linkStr)
}

func TestGenerateLinkConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        mtu: 9000
      - interface: eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	result := GenerateLinkConfig(m.Nodes[0].NetworkInterfaces[0])
	if result == nil {
		t.Fatal("expected link config, got nil")
	}

	if result.MetaName != "eth0" {
		t.Errorf("expected name=eth0, got %s", result.MetaName)
	}

	if result.LinkMTU != 9000 {
		t.Errorf("expected MTU 9000, got %d", result.LinkMTU)
	}

	result2 := GenerateLinkConfig(m.Nodes[0].NetworkInterfaces[1])
	if result2 != nil {
		t.Error("expected nil for interface without MTU")
	}
}

func TestGenerateLinkConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        mtu: 9000
      - interface: eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	linkBytes, err := GenerateLinkConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if linkBytes == nil {
		t.Fatal("expected link config bytes, got nil")
	}

	linkStr := string(linkBytes)
	if !bytes.Contains(linkBytes, []byte("kind: LinkConfig")) {
		t.Error("expected output to contain 'kind: LinkConfig'")
	}
	if !bytes.Contains(linkBytes, []byte("name: eth0")) {
		t.Error("expected output to contain 'name: eth0'")
	}
	if !bytes.Contains(linkBytes, []byte("mtu: 9000")) {
		t.Error("expected output to contain 'mtu: 9000'")
	}
	t.Logf("Link config output:\n%s", linkStr)
}

func TestGenerateVLANConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        vlans:
          - vlanId: 100
            addresses:
              - 192.168.100.1/24
            routes:
              - network: 10.0.0.0/8
                gateway: 192.168.100.254
            mtu: 1500`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	vlans := m.Nodes[0].NetworkInterfaces[0].DeviceVlans
	if len(vlans) == 0 {
		t.Fatal("expected VLANs, got none")
	}

	result := GenerateVLANConfig(m.Nodes[0].NetworkInterfaces[0], vlans[0])
	if result == nil {
		t.Fatal("expected VLAN config, got nil")
	}

	if result.MetaName != "eth0.100" {
		t.Errorf("expected name=eth0.100, got %s", result.MetaName)
	}

	if result.VLANIDConfig != 100 {
		t.Errorf("expected VLAN ID 100, got %d", result.VLANIDConfig)
	}

	if result.ParentLinkConfig != "eth0" {
		t.Errorf("expected parent=eth0, got %s", result.ParentLinkConfig)
	}

	if len(result.LinkAddresses) != 1 {
		t.Fatalf("expected 1 address, got %d", len(result.LinkAddresses))
	}

	expectedAddr := netip.MustParsePrefix("192.168.100.1/24")
	if result.LinkAddresses[0].AddressAddress != expectedAddr {
		t.Errorf("expected address %s, got %s", expectedAddr, result.LinkAddresses[0].AddressAddress)
	}

	if len(result.LinkRoutes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(result.LinkRoutes))
	}

	if result.LinkMTU != 1500 {
		t.Errorf("expected MTU 1500, got %d", result.LinkMTU)
	}
}

func TestGenerateVLANConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        vlans:
          - vlanId: 100
            addresses:
              - 192.168.100.1/24`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	vlanBytes, err := GenerateVLANConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if vlanBytes == nil {
		t.Fatal("expected VLAN config bytes, got nil")
	}

	vlanStr := string(vlanBytes)
	if !bytes.Contains(vlanBytes, []byte("kind: VLANConfig")) {
		t.Error("expected output to contain 'kind: VLANConfig'")
	}
	if !bytes.Contains(vlanBytes, []byte("name: eth0.100")) {
		t.Error("expected output to contain 'name: eth0.100'")
	}
	if !bytes.Contains(vlanBytes, []byte("vlanID: 100")) {
		t.Error("expected output to contain 'vlanID: 100'")
	}
	if !bytes.Contains(vlanBytes, []byte("parent: eth0")) {
		t.Error("expected output to contain 'parent: eth0'")
	}
	t.Logf("VLAN config output:\n%s", vlanStr)
}

func TestGenerateWireguardConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: wg0
        wireguard:
          privateKey: "iAmAPrivateKey="
          listenPort: 51820
          firewallMark: 51820
          peers:
            - publicKey: "iAmAPeerPublicKey="
              endpoint: "192.168.1.100:51820"
              persistentKeepaliveInterval: 10s
              allowedIPs:
                - 10.1.0.0/24`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	result := GenerateWireguardConfig(m.Nodes[0].NetworkInterfaces[0])
	if result == nil {
		t.Fatal("expected WireGuard config, got nil")
	}

	if result.MetaName != "wg0" {
		t.Errorf("expected name=wg0, got %s", result.MetaName)
	}

	if result.WireguardPrivateKey != "iAmAPrivateKey=" {
		t.Errorf("expected privateKey=iAmAPrivateKey=, got %s", result.WireguardPrivateKey)
	}

	if result.WireguardListenPort != 51820 {
		t.Errorf("expected listenPort=51820, got %d", result.WireguardListenPort)
	}

	if result.WireguardFirewallMark != 51820 {
		t.Errorf("expected firewallMark=51820, got %d", result.WireguardFirewallMark)
	}

	if len(result.WireguardPeers) != 1 {
		t.Fatalf("expected 1 peer, got %d", len(result.WireguardPeers))
	}

	peer := result.WireguardPeers[0]
	if peer.WireguardPublicKey != "iAmAPeerPublicKey=" {
		t.Errorf("expected peer publicKey=iAmAPeerPublicKey=, got %s", peer.WireguardPublicKey)
	}

	expectedEndpoint := netip.MustParseAddrPort("192.168.1.100:51820")
	if peer.WireguardEndpoint.AddrPort != expectedEndpoint {
		t.Errorf("expected endpoint=%s, got %s", expectedEndpoint, peer.WireguardEndpoint.AddrPort)
	}

	if len(peer.WireguardAllowedIPs) != 1 {
		t.Fatalf("expected 1 allowedIP, got %d", len(peer.WireguardAllowedIPs))
	}

	expectedAllowedIP := netip.MustParsePrefix("10.1.0.0/24")
	if peer.WireguardAllowedIPs[0].Prefix != expectedAllowedIP {
		t.Errorf("expected allowedIP=%s, got %s", expectedAllowedIP, peer.WireguardAllowedIPs[0].Prefix)
	}
}

func TestGenerateWireguardConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: wg0
        wireguard:
          privateKey: "iAmAPrivateKey="
          listenPort: 51820`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	wgBytes, err := GenerateWireguardConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if wgBytes == nil {
		t.Fatal("expected WireGuard config bytes, got nil")
	}

	wgStr := string(wgBytes)
	if !bytes.Contains(wgBytes, []byte("kind: WireguardConfig")) {
		t.Error("expected output to contain 'kind: WireguardConfig'")
	}
	if !bytes.Contains(wgBytes, []byte("name: wg0")) {
		t.Error("expected output to contain 'name: wg0'")
	}
	if !bytes.Contains(wgBytes, []byte("privateKey: iAmAPrivateKey=")) {
		t.Error("expected output to contain 'privateKey: iAmAPrivateKey='")
	}
	if !bytes.Contains(wgBytes, []byte("listenPort: 51820")) {
		t.Error("expected output to contain 'listenPort: 51820'")
	}
	t.Logf("WireGuard config output:\n%s", wgStr)
}

func TestGenerateBridgeConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: br0
        bridge:
          interfaces:
            - eth0
            - eth1
          stp:
            enabled: true`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	result := GenerateBridgeConfig(m.Nodes[0].NetworkInterfaces[0])
	if result == nil {
		t.Fatal("expected bridge config, got nil")
	}

	if result.MetaName != "br0" {
		t.Errorf("expected name=br0, got %s", result.MetaName)
	}

	if len(result.BridgeLinks) != 2 {
		t.Fatalf("expected 2 bridge links, got %d", len(result.BridgeLinks))
	}

	if result.BridgeLinks[0] != "eth0" {
		t.Errorf("expected first link=eth0, got %s", result.BridgeLinks[0])
	}

	if result.BridgeLinks[1] != "eth1" {
		t.Errorf("expected second link=eth1, got %s", result.BridgeLinks[1])
	}

	if result.BridgeSTP.BridgeSTPEnabled == nil || !*result.BridgeSTP.BridgeSTPEnabled {
		t.Error("expected STP enabled=true")
	}
}

func TestGenerateBridgeConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: br0
        bridge:
          interfaces:
            - eth0
            - eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	bridgeBytes, err := GenerateBridgeConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if bridgeBytes == nil {
		t.Fatal("expected bridge config bytes, got nil")
	}

	bridgeStr := string(bridgeBytes)
	if !bytes.Contains(bridgeBytes, []byte("kind: BridgeConfig")) {
		t.Error("expected output to contain 'kind: BridgeConfig'")
	}
	if !bytes.Contains(bridgeBytes, []byte("name: br0")) {
		t.Error("expected output to contain 'name: br0'")
	}
	if !bytes.Contains(bridgeBytes, []byte("- eth0")) {
		t.Error("expected output to contain '- eth0'")
	}
	if !bytes.Contains(bridgeBytes, []byte("- eth1")) {
		t.Error("expected output to contain '- eth1'")
	}
	t.Logf("Bridge config output:\n%s", bridgeStr)
}

func TestBondConfigWithAddressesAndRoutes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: bond0
        addresses:
          - 192.168.1.100/24
        routes:
          - network: 10.0.0.0/8
            gateway: 192.168.1.1
        mtu: 9000
        bond:
          interfaces:
            - eth0
            - eth1
          mode: 802.3ad`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	bondConfig := GenerateBondConfig(m.Nodes[0].NetworkInterfaces[0])
	if bondConfig == nil {
		t.Fatal("expected bond config, got nil")
	}

	if len(bondConfig.LinkAddresses) != 1 {
		t.Fatalf("expected 1 address in BondConfig, got %d", len(bondConfig.LinkAddresses))
	}
	expectedAddr := netip.MustParsePrefix("192.168.1.100/24")
	if bondConfig.LinkAddresses[0].AddressAddress != expectedAddr {
		t.Errorf("expected address %s, got %s", expectedAddr, bondConfig.LinkAddresses[0].AddressAddress)
	}

	if len(bondConfig.LinkRoutes) != 1 {
		t.Fatalf("expected 1 route in BondConfig, got %d", len(bondConfig.LinkRoutes))
	}
	expectedDest := netip.MustParsePrefix("10.0.0.0/8")
	if bondConfig.LinkRoutes[0].RouteDestination.Prefix != expectedDest {
		t.Errorf("expected destination %s, got %s", expectedDest, bondConfig.LinkRoutes[0].RouteDestination.Prefix)
	}

	if bondConfig.LinkMTU != 9000 {
		t.Errorf("expected MTU 9000, got %d", bondConfig.LinkMTU)
	}

	linkBytes, err := GenerateLinkConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}
	if linkBytes != nil {
		t.Error("expected no LinkConfig for bond interface, but got one")
	}
}

func TestWireguardConfigWithAddresses(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: wg0
        addresses:
          - 10.1.0.1/24
        mtu: 1420
        wireguard:
          privateKey: "iAmAPrivateKey="
          listenPort: 51820`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	wgConfig := GenerateWireguardConfig(m.Nodes[0].NetworkInterfaces[0])
	if wgConfig == nil {
		t.Fatal("expected wireguard config, got nil")
	}

	if len(wgConfig.LinkAddresses) != 1 {
		t.Fatalf("expected 1 address in WireguardConfig, got %d", len(wgConfig.LinkAddresses))
	}
	expectedAddr := netip.MustParsePrefix("10.1.0.1/24")
	if wgConfig.LinkAddresses[0].AddressAddress != expectedAddr {
		t.Errorf("expected address %s, got %s", expectedAddr, wgConfig.LinkAddresses[0].AddressAddress)
	}

	if wgConfig.LinkMTU != 1420 {
		t.Errorf("expected MTU 1420, got %d", wgConfig.LinkMTU)
	}

	linkBytes, err := GenerateLinkConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}
	if linkBytes != nil {
		t.Error("expected no LinkConfig for wireguard interface, but got one")
	}
}

func TestBridgeConfigWithMTU(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: br0
        mtu: 9000
        addresses:
          - 192.168.100.1/24
        bridge:
          interfaces:
            - eth0
            - eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	bridgeConfig := GenerateBridgeConfig(m.Nodes[0].NetworkInterfaces[0])
	if bridgeConfig == nil {
		t.Fatal("expected bridge config, got nil")
	}

	if bridgeConfig.LinkMTU != 9000 {
		t.Errorf("expected MTU 9000, got %d", bridgeConfig.LinkMTU)
	}

	if len(bridgeConfig.LinkAddresses) != 1 {
		t.Fatalf("expected 1 address in BridgeConfig, got %d", len(bridgeConfig.LinkAddresses))
	}

	linkBytes, err := GenerateLinkConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}
	if linkBytes != nil {
		t.Error("expected no LinkConfig for bridge interface, but got one")
	}
}

func TestHasSpecialConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: bond0
        bond:
          interfaces:
            - eth0
      - interface: wg0
        wireguard:
          privateKey: "test"
      - interface: br0
        bridge:
          interfaces:
            - eth1
      - interface: eth0
        vlans:
          - vlanId: 100
      - interface: eth2
        addresses:
          - 192.168.1.1/24`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if !hasSpecialConfig(m.Nodes[0].NetworkInterfaces[0]) {
		t.Error("expected bond interface to have special config")
	}

	if !hasSpecialConfig(m.Nodes[0].NetworkInterfaces[1]) {
		t.Error("expected wireguard interface to have special config")
	}

	if !hasSpecialConfig(m.Nodes[0].NetworkInterfaces[2]) {
		t.Error("expected bridge interface to have special config")
	}

	if !hasSpecialConfig(m.Nodes[0].NetworkInterfaces[3]) {
		t.Error("expected interface with VLAN to have special config")
	}

	if hasSpecialConfig(m.Nodes[0].NetworkInterfaces[4]) {
		t.Error("expected regular interface to NOT have special config")
	}
}

func TestGenerateLinkAliasConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        deviceSelector:
          hardwareAddr: "00:50:56:*"
      - interface: eth1
        deviceSelector:
          busPath: "0000:01:*"
          driver: virtio_net
      - interface: bond0
        deviceSelector:
          permanentAddr: "aa:bb:cc:dd:ee:ff"
          pciID: "8086:1533"`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	result1, err := GenerateLinkAliasConfig(m.Nodes[0].NetworkInterfaces[0], "eth0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result1 == nil {
		t.Fatal("expected link alias config, got nil")
	}
	if result1.MetaName != "eth0" {
		t.Errorf("expected name=eth0, got %s", result1.MetaName)
	}
	expr1 := result1.Selector.Match.String()
	t.Logf("CEL expression 1: %s", expr1)
	if !bytes.Contains([]byte(expr1), []byte(`glob(`)) || !bytes.Contains([]byte(expr1), []byte(`00:50:56:*`)) {
		t.Errorf("expected CEL expression to contain glob function with MAC pattern, got: %s", expr1)
	}

	result2, err := GenerateLinkAliasConfig(m.Nodes[0].NetworkInterfaces[1], "eth1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result2 == nil {
		t.Fatal("expected link alias config, got nil")
	}
	if result2.MetaName != "eth1" {
		t.Errorf("expected name=eth1, got %s", result2.MetaName)
	}
	expr2 := result2.Selector.Match.String()
	t.Logf("CEL expression 2: %s", expr2)
	if !bytes.Contains([]byte(expr2), []byte(`link.bus_path`)) {
		t.Errorf("expected expression to contain bus_path, got: %s", expr2)
	}
	if !bytes.Contains([]byte(expr2), []byte(`link.driver`)) {
		t.Errorf("expected expression to contain driver, got: %s", expr2)
	}
	if !bytes.Contains([]byte(expr2), []byte(`0000:01:*`)) {
		t.Errorf("expected expression to contain bus path pattern, got: %s", expr2)
	}
	if !bytes.Contains([]byte(expr2), []byte(`virtio_net`)) {
		t.Errorf("expected expression to contain driver name, got: %s", expr2)
	}

	result3, err := GenerateLinkAliasConfig(m.Nodes[0].NetworkInterfaces[2], "bond0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result3 == nil {
		t.Fatal("expected link alias config, got nil")
	}
	if result3.MetaName != "bond0" {
		t.Errorf("expected name=bond0, got %s", result3.MetaName)
	}
	expr3 := result3.Selector.Match.String()
	t.Logf("CEL expression 3: %s", expr3)
	if !bytes.Contains([]byte(expr3), []byte(`permanent_addr`)) {
		t.Errorf("expected expression to contain permanent_addr, got: %s", expr3)
	}
	if !bytes.Contains([]byte(expr3), []byte(`aa:bb:cc:dd:ee:ff`)) {
		t.Errorf("expected expression to contain MAC address, got: %s", expr3)
	}
	if !bytes.Contains([]byte(expr3), []byte(`link.pciid`)) {
		t.Errorf("expected expression to contain pciid, got: %s", expr3)
	}
	if !bytes.Contains([]byte(expr3), []byte(`8086:1533`)) {
		t.Errorf("expected expression to contain PCI ID, got: %s", expr3)
	}
}

func TestGenerateLinkAliasConfigWithoutInterfaceName(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - deviceSelector:
          hardwareAddr: "00:11:22:33:44:55"
      - deviceSelector:
          driver: "e1000"
        bond:
          mode: 802.3ad
          interfaces: []
      - deviceSelector:
          busPath: "0000:03:*"
        bridge:
          interfaces: []`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	// Eth interface
	result1, err := GenerateLinkAliasConfig(m.Nodes[0].NetworkInterfaces[0], "eth0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result1 == nil {
		t.Fatal("expected link alias config, got nil")
	}
	if result1.MetaName != "eth0" {
		t.Errorf("expected synthetic name=eth0, got %s", result1.MetaName)
	}

	// Bond interface
	result2, err := GenerateLinkAliasConfig(m.Nodes[0].NetworkInterfaces[1], "bond0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result2 == nil {
		t.Fatal("expected link alias config, got nil")
	}
	if result2.MetaName != "bond0" {
		t.Errorf("expected synthetic name=bond0, got %s", result2.MetaName)
	}

	// Bridge interface
	result3, err := GenerateLinkAliasConfig(m.Nodes[0].NetworkInterfaces[2], "br0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result3 == nil {
		t.Fatal("expected link alias config, got nil")
	}
	if result3.MetaName != "br0" {
		t.Errorf("expected synthetic name=br0, got %s", result3.MetaName)
	}
}

func TestGenerateLinkAliasConfigAvoidsDuplicateNames(t *testing.T) {
	data := []byte(`nodes:
  - hostname: test-node
    networkInterfaces:
      - interface: eth0
        addresses:
          - 192.168.1.10/24
      - deviceSelector:
          hardwareAddr: "00:30:93:*"
        addresses:
          - 10.0.0.10/24
      - deviceSelector:
          hardwareAddr: "00:50:56:*"
        addresses:
          - 172.16.0.10/24`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	aliasBytes, err := GenerateLinkAliasConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	aliasConfigs := string(aliasBytes)
	t.Logf("Generated LinkAlias configs:\n%s", aliasConfigs)

	linkAliasCount := bytes.Count(aliasBytes, []byte("kind: LinkAliasConfig"))
	if linkAliasCount != 2 {
		t.Errorf("expected 2 LinkAliasConfig documents (for 2nd and 3rd interfaces), got %d", linkAliasCount)
	}

	lines := bytes.Split(aliasBytes, []byte("\n"))
	for _, line := range lines {
		if bytes.Contains(line, []byte("name: eth0")) {
			t.Error("found 'name: eth0' in generated LinkAlias - should not generate alias for explicitly named interface")
		}
	}

	if !bytes.Contains(aliasBytes, []byte("name: ethSel0")) {
		t.Error("expected second interface (first with deviceSelector) to be named ethSel0")
	}
	if !bytes.Contains(aliasBytes, []byte("name: ethSel1")) {
		t.Error("expected third interface (second with deviceSelector) to be named ethSel1")
	}
}

func TestDeviceSelectorUpdatesDeviceInterface(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: eth0
        addresses:
          - 192.168.1.10/24
      - deviceSelector:
          hardwareAddr: "00:50:56:*"
        addresses:
          - 10.0.0.10/24`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if len(m.Nodes) == 0 || len(m.Nodes[0].NetworkInterfaces) != 2 {
		t.Fatal("failed to parse test data")
	}

	devices := m.Nodes[0].NetworkInterfaces

	_, err := GenerateLinkAliasConfigBytes(devices)
	if err != nil {
		t.Fatalf("failed to generate link alias config: %v", err)
	}

	if devices[1].DeviceInterface == "" {
		t.Error("expected Device.DeviceInterface to be updated with synthetic name, but it's still empty")
	}
	if devices[1].DeviceInterface != "ethSel0" {
		t.Errorf("expected Device.DeviceInterface to be 'ethSel0', got '%s'", devices[1].DeviceInterface)
	}

	linkBytes, err := GenerateLinkConfigBytes(devices)
	if err != nil {
		t.Fatalf("failed to generate link config: %v", err)
	}

	if !bytes.Contains(linkBytes, []byte("name: eth0")) {
		t.Error("expected LinkConfig for eth0")
	}
	if !bytes.Contains(linkBytes, []byte("name: ethSel0")) {
		t.Error("expected LinkConfig for ethSel0 (synthetic name)")
	}

	emptyNameCount := bytes.Count(linkBytes, []byte("name: \"\""))
	if emptyNameCount > 0 {
		t.Errorf("found %d LinkConfig with empty name, expected 0", emptyNameCount)
	}
}

func TestCELExpressionEvaluation(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		mockLink   *networktypes.LinkStatusSpec
		expected   bool
	}{
		{
			name:       "driver match",
			expression: `glob("e1000", link.driver)`,
			mockLink: &networktypes.LinkStatusSpec{
				Driver: "e1000",
			},
			expected: true,
		},
		{
			name:       "driver no match",
			expression: `glob("e1000", link.driver)`,
			mockLink: &networktypes.LinkStatusSpec{
				Driver: "virtio_net",
			},
			expected: false,
		},
		{
			name:       "MAC address wildcard match",
			expression: `glob("00:50:56:*", mac(link.hardware_addr))`,
			mockLink: &networktypes.LinkStatusSpec{
				HardwareAddr: []byte{0x00, 0x50, 0x56, 0x12, 0x34, 0x56},
			},
			expected: true,
		},
		{
			name:       "MAC address no match",
			expression: `glob("aa:bb:cc:*", mac(link.hardware_addr))`,
			mockLink: &networktypes.LinkStatusSpec{
				HardwareAddr: []byte{0x00, 0x50, 0x56, 0x12, 0x34, 0x56},
			},
			expected: false,
		},
		{
			name:       "compound expression - both match",
			expression: `glob("e1000", link.driver) && glob("00:50:56:*", mac(link.hardware_addr))`,
			mockLink: &networktypes.LinkStatusSpec{
				Driver:       "e1000",
				HardwareAddr: []byte{0x00, 0x50, 0x56, 0x12, 0x34, 0x56},
			},
			expected: true,
		},
		{
			name:       "compound expression - one fails",
			expression: `glob("virtio_net", link.driver) && glob("00:50:56:*", mac(link.hardware_addr))`,
			mockLink: &networktypes.LinkStatusSpec{
				Driver:       "e1000",
				HardwareAddr: []byte{0x00, 0x50, 0x56, 0x12, 0x34, 0x56},
			},
			expected: false,
		},
		{
			name:       "bus path match",
			expression: `glob("0000:01:*", link.bus_path)`,
			mockLink: &networktypes.LinkStatusSpec{
				BusPath: "0000:01:00.0",
			},
			expected: true,
		},
		{
			name:       "PCI ID match",
			expression: `glob("8086:1533", link.pciid)`,
			mockLink: &networktypes.LinkStatusSpec{
				Pciid: "8086:1533",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			celExpr, err := cel.ParseBooleanExpression(tt.expression, celenv.LinkLocator())
			if err != nil {
				t.Fatalf("failed to parse CEL expression: %v", err)
			}

			result, err := celExpr.EvalBool(celenv.LinkLocator(), map[string]any{
				"link": tt.mockLink,
			})
			if err != nil {
				t.Fatalf("failed to evaluate CEL expression: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBuildDeviceSelectorCELExpression_Physical(t *testing.T) {
	sel := &v1alpha1.NetworkDeviceSelector{NetworkDevicePhysical: pointer.To(true)}

	expr, err := buildDeviceSelectorCELExpression(sel)
	if err != nil {
		t.Fatalf("failed to build CEL expression: %v", err)
	}
	if expr == "" {
		t.Fatalf("expected non-empty expression for physical selector")
	}

	celExpr, err := cel.ParseBooleanExpression(expr, celenv.LinkLocator())
	if err != nil {
		t.Fatalf("failed to parse CEL expression: %v", err)
	}

	result, err := celExpr.EvalBool(celenv.LinkLocator(), map[string]any{
		"link": &networktypes.LinkStatusSpec{Type: netenums.NethelpersLinkType_LINK_ETHER},
	})
	if err != nil {
		t.Fatalf("failed to evaluate CEL expression: %v", err)
	}

	if !result {
		t.Fatalf("expected physical selector to match an ether link")
	}

	sel2 := &v1alpha1.NetworkDeviceSelector{NetworkDevicePhysical: pointer.To(false)}
	expr2, err := buildDeviceSelectorCELExpression(sel2)
	if err != nil {
		t.Fatalf("failed to build CEL expression: %v", err)
	}
	celExpr2, err := cel.ParseBooleanExpression(expr2, celenv.LinkLocator())
	if err != nil {
		t.Fatalf("failed to parse CEL expression: %v", err)
	}

	result2, err := celExpr2.EvalBool(celenv.LinkLocator(), map[string]any{
		"link": &networktypes.LinkStatusSpec{Type: netenums.NethelpersLinkType_LINK_ETHER},
	})
	if err != nil {
		t.Fatalf("failed to evaluate CEL expression: %v", err)
	}
	if result2 {
		t.Fatalf("expected physical=false selector to NOT match an ether link")
	}
}

func TestGenerateLinkAliasConfigBytes(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: net0
        deviceSelector:
          hardwareAddr: "00:50:56:12:34:56"
          driver: e1000
        dhcp: true
      - interface: eth1`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	aliasBytes, err := GenerateLinkAliasConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if aliasBytes == nil {
		t.Fatal("expected link alias config bytes, got nil")
	}

	aliasStr := string(aliasBytes)
	if !bytes.Contains(aliasBytes, []byte("kind: LinkAliasConfig")) {
		t.Error("expected output to contain 'kind: LinkAliasConfig'")
	}
	if !bytes.Contains(aliasBytes, []byte("name: net0")) {
		t.Error("expected output to contain 'name: net0'")
	}
	if !bytes.Contains(aliasBytes, []byte("match:")) {
		t.Error("expected output to contain 'match:'")
	}
	t.Logf("Link alias config output:\n%s", aliasStr)
}

func TestGenerateBondMemberAliasConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    networkInterfaces:
      - interface: bond0
        bond:
          deviceSelectors:
            - hardwareAddr: "58:47:ca:76:9c:*"
              driver: i40e
            - hardwareAddr: "58:47:ca:76:9d:*"
              driver: i40e
          mode: 802.3ad`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	bondConfig := GenerateBondConfig(m.Nodes[0].NetworkInterfaces[0])
	if bondConfig == nil {
		t.Fatal("expected bond config, got nil")
	}

	if len(bondConfig.BondLinks) != 2 {
		t.Fatalf("expected 2 bond links, got %d", len(bondConfig.BondLinks))
	}

	if bondConfig.BondLinks[0] != "bond0-m0" {
		t.Errorf("expected first link=bond0-m0, got %s", bondConfig.BondLinks[0])
	}

	if bondConfig.BondLinks[1] != "bond0-m1" {
		t.Errorf("expected second link=bond0-m1, got %s", bondConfig.BondLinks[1])
	}

	aliasBytes, err := GenerateBondMemberAliasConfigBytes(m.Nodes[0].NetworkInterfaces)
	if err != nil {
		t.Fatal(err)
	}

	if aliasBytes == nil {
		t.Fatal("expected bond member alias config bytes, got nil")
	}

	aliasStr := string(aliasBytes)
	t.Logf("Bond member alias config output:\n%s", aliasStr)

	if !bytes.Contains(aliasBytes, []byte("name: bond0-m0")) {
		t.Error("expected output to contain 'name: bond0-m0'")
	}
	if !bytes.Contains(aliasBytes, []byte("name: bond0-m1")) {
		t.Error("expected output to contain 'name: bond0-m1'")
	}
	if !bytes.Contains(aliasBytes, []byte("58:47:ca:76:9c:*")) {
		t.Error("expected output to contain first MAC pattern")
	}
	if !bytes.Contains(aliasBytes, []byte("58:47:ca:76:9d:*")) {
		t.Error("expected output to contain second MAC pattern")
	}
}
