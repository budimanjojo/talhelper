package validate

import (
	"regexp"
	"strings"

	"github.com/gookit/validate"
	"github.com/talos-systems/net"
)

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

func (c Config) IsVersion(version string) bool {
	if match, _ := regexp.MatchString(`^v?(\d+\.)(\d+\.)(\d+)$`, version); match {
		return true
	}
	return false
}

func (c Config) IsCNIName(cni string) bool {
	if match, _ := regexp.MatchString(`^none$|^flannel$|^custom$`, cni); match {
		return true
	}
	return false
}

func (c Config) IsCIDRList(nets []string) bool {
	for _, net := range nets {
		if !validate.IsCIDR(net) {
			return false
		}
	}
	return true
}

func (c Config) IsIPList(ips []string) bool {
	for _, ip := range ips {
		if !validate.IsIP(ip) {
			return false
		}
	}
	return true
}

func (c Config) IsURLList(urls []string) bool {
	for _, url := range urls {
		if !validate.IsURL(url) {
			return false
		}
	}
	return true
}

func (c Config) IsTalosEndpoint(ep string) bool {
	if err := net.ValidateEndpointURI(ep); err != nil {
		return false
	}
	return true
}

func (c Config) IsDomain(domain string) bool {
	if domain == "" || len(domain)-strings.Count(domain, ".") > 255 {
		return false
	}
	return regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`).MatchString(domain)
}
