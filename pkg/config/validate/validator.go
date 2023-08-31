package validate

import (
	"regexp"
	"strings"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/gookit/validate"
	"github.com/siderolabs/net"
	"github.com/siderolabs/talos/pkg/machinery/api/machine"
	"github.com/siderolabs/talos/pkg/machinery/compatibility"
	"golang.org/x/mod/semver"
)

// IsRFC6902List returns true if `input` is list of RFC6902 JSON patch.
func (c Config) IsRFC6902List(input []map[string]interface{}) bool {
	for _, v := range input {
		if _, ok := v["path"]; ok {
			if val, ok := v["op"]; ok {
				switch val {
				case "add":
					if _, ok := v["value"]; ok {
						continue
					}
					return false
				case "remove":
					continue
				default:
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

// IsSemVer returns true if `version` is a valid semantic version.
func (c Config) IsSemVer(version string) bool {
	stripped := strings.TrimPrefix(version, "v")
	re := `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	if match, _ := regexp.MatchString(re, stripped); match {
		return true
	}
	return false
}

// IsSupportedTalosVersion returns true if `version` is supported Talos version.
func (c Config) IsSupportedTalosVersion(version string) bool {
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	majorMinor := semver.MajorMinor(version)

	switch majorMinor {
	case "v1.2":
		return true
	case "v1.3":
		return true
	case "v1.4":
		return true
	case "v1.5":
		return true
	case "v1.6":
		return true
	default:
		return false
	}
}

// IsSupportedK8sVersion returns true if Kubernetes `version` is supported by
// `c.TalosVersion`.
func (c Config) IsSupportedK8sVersion(version string) bool {
	var talosVersionInfo *machine.VersionInfo

	// stop here if `c.TalosVersion` is not right
	if !c.IsSemVer(c.TalosVersion) {
		return false
	}

	if c.TalosVersion == "" {
		talosVersionInfo = &machine.VersionInfo{
			Tag: config.LatestTalosVersion,
		}
	} else {
		talosVersionInfo = &machine.VersionInfo{
			Tag: c.TalosVersion,
		}
	}

	talosVersion, _ := compatibility.ParseTalosVersion(talosVersionInfo)

	kubernetesVersion, err := compatibility.ParseKubernetesVersion(strings.TrimPrefix(version, "v"))
	if err != nil {
		return false
	}

	if err := kubernetesVersion.SupportedWith(talosVersion); err != nil {
		return false
	}
	return true
}

// IsCNIName returns true if `cni` is a supported Talos CNI name.
func (c Config) IsCNIName(cni string) bool {
	if match, _ := regexp.MatchString(`^none$|^flannel$|^custom$`, cni); match {
		return true
	}
	return false
}

// IsValidFileOperation returns true if `op` is a correct `machineFile` operation
func (c Config) IsValidFileOperation(op string) bool {
	if match, _ := regexp.MatchString(`^create$|^append$|^overwrite$`, op); match {
		return true
	}
	return false
}

// IsCIDRList returns true if `nets` is list of CIDR addresses.
func (c Config) IsCIDRList(nets []string) bool {
	for _, net := range nets {
		if !validate.IsCIDR(net) {
			return false
		}
	}
	return true
}

// IsIPList returns true if `ips` is list of IP addresses.
func (c Config) IsIPList(ips []string) bool {
	for _, ip := range ips {
		if !validate.IsIP(ip) {
			return false
		}
	}
	return true
}

// IsURLList returns true if `urls` is list of URLs.
func (c Config) IsURLList(urls []string) bool {
	for _, url := range urls {
		if !validate.IsURL(url) {
			return false
		}
	}
	return true
}

// IsTalosEndpoint returns true if `ep` is a valid Talos endpoint.
func (c Config) IsTalosEndpoint(ep string) bool {
	if err := net.ValidateEndpointURI(ep); err != nil {
		return false
	}
	return true
}

// IsDomain returns true if `domain` is a valid domain.
func (c Config) IsDomain(domain string) bool {
	if domain == "" || len(domain)-strings.Count(domain, ".") > 255 {
		return false
	}
	return regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`).MatchString(domain)
}

// IsDomainOrIP returns true if `domainIP` is either valid domain or IP.
func (c Config) IsDomainOrIP(domainIP string) bool {
	if c.IsDomain(domainIP) || validate.IsIP(domainIP) {
		return true
	}

	return false
}

// IsValidNetworkInterfaces returns true if `ifaces` is list of network interfaces.
func (c Config) IsValidNetworkInterfaces(ifaces []*NetworkInterface) bool {
	for _, iface := range ifaces {
		if iface.Interface == "" && iface.DeviceSelector == nil {
			return false
		} else if iface.Interface != "" && iface.DeviceSelector != nil {
			return false
		}
	}
	return true
}
