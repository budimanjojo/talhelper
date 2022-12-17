package validate

import (
	"time"

	"github.com/gookit/validate"
)

type Config struct {
	ClusterName                    string   `validate:"required"`
	TalosVersion                   string   `validate:"isVersion"`
	KubernetesVersion              string   `validate:"isVersion"`
	Endpoint                       string   `validate:"isTalosEndpoint"`
	Domain                         string   `validate:"isDomain"`
	AllowSchedulingOnMasters       string   `validate:"isBool"`
	AllowSchedulingOnControlPlanes string   `validate:"isBool"`
	ClusterPodNets                 []string `validate:"isCIDRList"`
	ClusterSvcNets                 []string `validate:"isCIDRList"`
	CniConfig                      *CNIConfig
	Nodes                          []*Node
	ControlPlane                   struct {
		ConfigPatches []map[string]interface{} `validate:"isRFC6902List"`
	}
	Worker struct {
		ConfigPatches []map[string]interface{} `validate:"isRFC6902List"`
	}
}

type CNIConfig struct {
	Name string   `validate:"isCNIName|requiredWith:CniConfig"`
	Urls []string `validate:"isURLList|requiredIf:CniConfig.Name,custom"`
}

type Node struct {
	Hostname            string                   `validate:"required"`
	IPAddress           string                   `validate:"required|isDomainOrIP"`
	ControlPlane        string                   `validate:"isBool"`
	InstallDisk         string                   `validate:"requiredWithout:Nodes.InstallDiskSelector"`
	DisableSearchDomain string                   `validate:"isBool"`
	Nameservers         []string                 `validate:"isIPList"`
	ConfigPatches       []map[string]interface{} `validate:"isRFC6902List"`
	NetworkInterfaces   []*NetworkInterface
	InstallDiskSelector *InstallDiskSelector
	KernelModules       []*KernelModule
}

type KernelModule struct {
	Name       string
	Parameters []string
}

type InstallDiskSelector struct {
	Size     string
	Name     string
	Model    string
	Modalias string
	UUID     string
	WWID     string
	Type     string
	BusPath  string
}

type NetworkInterface struct {
	Interface   string   `validate:"required_with:Nodes.NetworkInterfaces"`
	Addresses   []string `validate:"isCIDRList"`
	Routes      []Route
	Bond        *Bond
	Vlans       []*Vlan
	Mtu         string `validate:"isIntString"`
	Dhcp        string `validate:"isBool"`
	Ignore      string `validate:"isBool"`
	Dummy       string `validate:"isBool"`
	DhcpOptions *DhcpOption
	Wireguard   *Wireguard
	Vip         *Vip
}

type Route struct {
	Network string `validate:"isCIDR"`
	Gateway string `validate:"isIP"`
	Source  string `validate:"isIP"`
	Metric  string `validate:"isUint"`
}

type Bond struct {
	Interfaces      []string `validate:"requiredWith:Nodes.NetworkInterfaces.Bond"`
	ArpIPTarget     []string
	Mode            string
	XmitHashPolicy  string
	LacpRate        string
	AdActorSystem   string
	ArpValidate     string
	Primary         string
	PrimaryReselect string
	FailOverMac     string
	AdSelect        string
	MiiMon          string `validate:"isUint"`
	Updelay         string `validate:"isUint"`
	Downdelay       string `validate:"isUint"`
	ArpInterval     string `validate:"isUint"`
	ResendIgmp      string `validate:"isUint"`
	MinLinks        string `validate:"isUint"`
	LpInterval      string `validate:"isUint"`
	PacketsPerSlave string `validate:"isUint"`
	NumPeerNotif    string `validate:"isUint"`
	TlbDynamicLb    string `validate:"isUint"`
	AllSlavesActive string `validate:"isUint"`
	UseCarrier      string `validate:"isBool"`
	AdActorSysPrio  string `validate:"isUint"`
	AdUserPortKey   string `validate:"isUint"`
	PeerNotifyDelay string `validate:"isUint"`
}

type Vip struct {
	Ip           string `validate:"isIP"`
	EquinixMetal struct {
		ApiToken string
	}
	Hcloud struct {
		ApiToken string
	}
}

type Vlan struct {
	Addresses []string `validate:"isCIDRList"`
	Routes    []*Route
	Dhcp      string `validate:"isBool"`
	VlanId    string `validate:"isUint"`
	Mtu       string `validate:"isUint"`
	Vip       *Vip
}

type DhcpOption struct {
	RouteMetric string `validate:"isUint"`
	Ipv4        string `validate:"isBool"`
	Ipv6        string `validate:"isBool"`
}

type Wireguard struct {
	PrivateKey   string
	ListenPort   string `validate:"isIntString"`
	FirewallMark string `validate:"isIntString"`
	Peers        []*Peer
}

type Peer struct {
	PublicKey                   string
	Endpoint                    string
	PersistentKeepAliveInterval time.Duration
	AllowedIPs                  []string
}

func (c Config) Messages() map[string]string {
	return validate.MS{
		"required":                  "{field} is required",
		"CniConfig.Urls.requiredIf": "{field} is required when {args0} is \"custom\"",
		"isRFC6902List":             "{field} doesn't look like list of RFC 6902 JSON patches",
		"isCIDR":                    "{field} doesn't look like CIDR notation",
		"isCIDRList":                "{field} doesn't look like list of CIDR notations",
		"isIP":                      "{field} doesn't look like IP address",
		"isIPList":                  "{field} doesn't look like list of IP addresses",
		"isURLList":                 "{field} doesn't look like list of URLs",
		"isTalosEndpoint":           "{field} is not a valid Talos endpoint",
		"isDomain":                  "{field} is not a valid domain",
		"isCNIName":                 "{field} is not a valid CNI name (none,flannel,custom)",
		"isBool":                    "{field} is not a valid boolean (true or false)",
		"isInt":                     "{field} is not a valid integer",
		"isUint":                    "{field} is not a valid unsigned integer",
		"isDomainOrIP":              "{field} is not a valid domain or IP address",
	}
}
