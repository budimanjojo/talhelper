package talos

import (
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
