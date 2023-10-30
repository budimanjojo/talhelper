package config

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/gookit/validate"
	"github.com/hashicorp/go-multierror"
	"github.com/siderolabs/net"
	"github.com/siderolabs/talos/pkg/machinery/api/machine"
	"github.com/siderolabs/talos/pkg/machinery/compatibility"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/siderolabs/talos/pkg/machinery/labels"
	"golang.org/x/mod/semver"
)

func checkRequiredCfg(c TalhelperConfig, result *Errors) *Errors {
	if c.ClusterName == "" {
		e := &Error{
			Kind:  "ClusterNameRequired",
			Field: getFieldYamlTag(c, "ClusterName"),
		}
		e.Message = formatError(multierror.Append(fmt.Errorf("%q is required to be not empty", e.Field)))
		result.Append(e)
	}

	if c.Endpoint == "" {
		e := &Error{
			Kind:  "EndpointRequired",
			Field: getFieldYamlTag(c, "Endpoint"),
		}
		e.Message = formatError(multierror.Append(fmt.Errorf("%q is required to be not empty", e.Field)))
		result.Append(e)
	}

	if len(c.Nodes) == 0 {
		e := &Error{
			Kind:  "NodesRequired",
			Field: getFieldYamlTag(c, "Nodes"),
		}
		e.Message = formatError(multierror.Append(fmt.Errorf("%q is required to be not empty", e.Field)))
		result.Append(e)
	}

	return result
}

func checkSupportedTalosVersion(c TalhelperConfig, result *Errors) *Errors {
	if c.TalosVersion != "" {
		if !strings.HasPrefix(c.TalosVersion, "v") {
			c.TalosVersion = "v" + c.TalosVersion
		}
		majorMinor := semver.MajorMinor(c.TalosVersion)
		switch majorMinor {
		case "v1.2":
			return result
		case "v1.3":
			return result
		case "v1.4":
			return result
		case "v1.5":
			return result
		case "v1.6":
			return result
		default:
			return result.Append(&Error{
				Kind:    "InvalidTalosVersion",
				Field:   getFieldYamlTag(c, "TalosVersion"),
				Message: formatError(multierror.Append(fmt.Errorf("%q is not a supported Talos version", c.TalosVersion))),
			})
		}
	}
	return result
}

func checkSupportedK8sVersion(c TalhelperConfig, result *Errors) *Errors {
	if c.KubernetesVersion != "" {
		var (
			messages         *multierror.Error
			talosVersionInfo *machine.VersionInfo
		)

		// stop here if `c.TalosVersion` is not right
		if c.TalosVersion != "" && !isSemVer(c.TalosVersion) {
			return result.Append(&Error{
				Kind:    "InvalidKubernetesVersion",
				Field:   getFieldYamlTag(c, "KubernetesVersion"),
				Message: formatError(multierror.Append(fmt.Errorf("fix the issue on %q field", getFieldYamlTag(c, "TalosVersion")))),
			})
		}

		if c.TalosVersion == "" {
			talosVersionInfo = &machine.VersionInfo{
				Tag: LatestTalosVersion,
			}
		} else {
			talosVersionInfo = &machine.VersionInfo{
				Tag: c.TalosVersion,
			}
		}

		talosVersion, err := compatibility.ParseTalosVersion(talosVersionInfo)
		if err != nil {
			messages = multierror.Append(messages, err)
		}

		kubernetesVersion, err := compatibility.ParseKubernetesVersion(strings.TrimPrefix(c.KubernetesVersion, "v"))
		if err != nil {
			messages = multierror.Append(messages, err)
		}

		if err := kubernetesVersion.SupportedWith(talosVersion); err != nil {
			messages = multierror.Append(messages, err)
		}

		if messages.ErrorOrNil() != nil {
			return result.Append(&Error{
				Kind:    "InvalidKubernetesVersion",
				Field:   getFieldYamlTag(c, "KubernetesVersion"),
				Message: formatError(messages),
			})
		}
		return result
	}
	return result
}

func checkTalosEndpoint(c TalhelperConfig, result *Errors) *Errors {
	if c.Endpoint != "" {
		var messages *multierror.Error

		if err := net.ValidateEndpointURI(c.Endpoint); err != nil {
			messages = multierror.Append(messages, err)
		}

		if messages.ErrorOrNil() != nil {
			return result.Append(&Error{
				Kind:    "InvalidTalosEndpoint",
				Field:   getFieldYamlTag(c, "Endpoint"),
				Message: formatError(messages),
			})
		}
	}
	return result
}

func checkDomain(c TalhelperConfig, result *Errors) *Errors {
	if c.Domain != "" {
		if !isDomain(c.Domain) {
			return result.Append(&Error{
				Kind:    "InvalidDomain",
				Field:   getFieldYamlTag(c, "Domain"),
				Message: formatError(multierror.Append(fmt.Errorf("%q is not a valid domain", c.Domain))),
			})
		}
	}
	return result
}

func checkClusterNets(c TalhelperConfig, result *Errors) *Errors {
	if len(c.ClusterPodNets) > 0 {
		if !isCIDRList(c.ClusterPodNets) {
			result = result.Append(&Error{
				Kind:    "InvalidClusterPodNets",
				Field:   getFieldYamlTag(c, "ClusterPodNets"),
				Message: formatError(multierror.Append(fmt.Errorf("%q doesn't look like list of CIDR notations", c.ClusterPodNets))),
			})
		}
	}

	if len(c.ClusterSvcNets) > 0 {
		if !isCIDRList(c.ClusterSvcNets) {
			result = result.Append(&Error{
				Kind:    "InvalidClusterSvcNets",
				Field:   getFieldYamlTag(c, "ClusterSvcNets"),
				Message: formatError(multierror.Append(fmt.Errorf("%q doesn't look like list of CIDR notations", c.ClusterSvcNets))),
			})
		}
	}
	return result
}

func checkCNIConfig(c TalhelperConfig, result *Errors) *Errors {
	if c.CNIConfig.Name != "" {
		var messages *multierror.Error

		switch c.CNIConfig.Name {
		case "flannel":
			fallthrough
		case "none":
			if len(c.CNIConfig.Urls) != 0 {
				messages = multierror.Append(messages, fmt.Errorf(`"Urls" field should be empty for %q CNI`, c.CNIConfig.Name))
			}
		case "custom":
			if len(c.CNIConfig.Urls) == 0 {
				messages = multierror.Append(messages, fmt.Errorf(`"Urls" field should not be empty for %q CNI`, c.CNIConfig.Name))
			}
			for _, url := range c.CNIConfig.Urls {
				if !validate.IsURL(url) {
					messages = multierror.Append(messages, fmt.Errorf(`%q in "Urls" field is not a valid url`, url))
				}
			}
		default:
			messages = multierror.Append(messages, fmt.Errorf("%q is not a valid CNI name (none,flannel,custom)", c.CNIConfig.Name))
		}

		if messages.ErrorOrNil() != nil {
			return result.Append(&Error{
				Kind:    "InvalidCNIConfig",
				Field:   getFieldYamlTag(c, "CNIConfig"),
				Message: formatError(messages),
			})
		}
	}
	return result
}

func checkControlPlane(c TalhelperConfig, result *Errors) *Errors {
	if len(c.ControlPlane.ConfigPatches) > 0 {
		if !isRFC6902List(c.ControlPlane.ConfigPatches) {
			result = result.Append(&Error{
				Kind:    "InvalidControlPlaneConfigPatches",
				Field:   getFieldYamlTag(c, "ControlPlane.ConfigPatches"),
				Message: formatError(multierror.Append(fmt.Errorf("doesn't look like list of RFC6902 JSON patches"))),
			})
		}
	}
	return result
}

func checkWorker(c TalhelperConfig, result *Errors) *Errors {
	if len(c.Worker.ConfigPatches) > 0 {
		if !isRFC6902List(c.Worker.ConfigPatches) {
			result = result.Append(&Error{
				Kind:    "InvalidWorkerConfigPatches",
				Field:   getFieldYamlTag(c, "Worker.ConfigPatches"),
				Message: formatError(multierror.Append(fmt.Errorf("doesn't look like list of RFC6902 JSON patches"))),
			})
		}
	}
	return result
}

func checkNodeRequiredCfg(node Node, idx int, result *Errors) *Errors {
	if node.Hostname == "" {
		e := &Error{
			Kind:  "NodeHostnameRequired",
			Field: getNodeFieldYamlTag(node, idx, "Hostname"),
		}
		e.Message = formatError(multierror.Append(fmt.Errorf("%q is required to be not empty", e.Field)))
		result = result.Append(e)
	}

	if node.IPAddress == "" {
		e := &Error{
			Kind:  "NodeIPAddressRequired",
			Field: getNodeFieldYamlTag(node, idx, "IPAddress"),
		}
		e.Message = formatError(multierror.Append(fmt.Errorf("%q is required to be not empty", e.Field)))
		result = result.Append(e)
	}

	if node.InstallDisk == "" && node.InstallDiskSelector == nil {
		e := &Error{
			Kind:  "NodeInstallRequired",
			Field: getNodeFieldYamlTag(node, idx, "InstallDisk"),
		}
		e.Message = formatError(multierror.Append(fmt.Errorf("%q is required to be not empty", e.Field)))
		result = result.Append(e)
	}

	return result
}

func checkNodeIPAddress(node Node, idx int, result *Errors) *Errors {
	if node.IPAddress != "" {
		if !isDomainOrIP(node.IPAddress) {
			return result.Append(&Error{
				Kind:    "InvalidNodeIPAddress",
				Field:   getNodeFieldYamlTag(node, idx, "IPAddress"),
				Message: formatError(multierror.Append(fmt.Errorf("%q is not a valid domain or IP address", node.IPAddress))),
			})
		}
	}
	return result
}

func checkNodeLabels(node Node, idx int, result *Errors) *Errors {
	if node.NodeLabels != nil {
		var messages *multierror.Error
		if err := labels.Validate(node.NodeLabels); err != nil {
			return result.Append(&Error{
				Kind:    "InvalidNodeLabels",
				Field:   getNodeFieldYamlTag(node, idx, "NodeLabels"),
				Message: formatError(multierror.Append(messages, err)),
			})
		}
	}
	return result
}

func checkNodeMachineDisks(node Node, idx int, result *Errors) *Errors {
	if node.MachineDisks != nil {
		var messages *multierror.Error

		for _, disk := range node.MachineDisks {
			for i, pt := range disk.DiskPartitions {
				if pt.DiskSize == 0 && i != len(disk.DiskPartitions)-1 {
					messages = multierror.Append(messages, fmt.Errorf("partition %q for disk %q is set to occupy full disk, but it's not the last partition in the list", pt.DiskMountPoint, disk.Device()))
				}
			}
		}

		if messages.ErrorOrNil() != nil {
			return result.Append(&Error{
				Kind:    "InvalidMachineDisks",
				Field:   getNodeFieldYamlTag(node, idx, "MachineDisks"),
				Message: formatError(messages),
			})
		}
	}
	return result
}

func checkNodeMachineFiles(node Node, idx int, result *Errors) *Errors {
	if node.MachineFiles != nil {
		var messages *multierror.Error
		pattern := `^create$|^append$|^overwrite$`
		re := regexp.MustCompile(pattern)

		for _, file := range node.MachineFiles {
			if !re.MatchString(file.FileOp) {
				messages = multierror.Append(messages, fmt.Errorf("%q is not a valid operation name (create,append,overwrite)", file.Op()))
			}
			if !validate.IsUnixPath(file.Path()) {
				messages = multierror.Append(messages, fmt.Errorf("%q is not a valid Unix file path", file.Path()))
			}
		}

		if messages.ErrorOrNil() != nil {
			return result.Append(&Error{
				Kind:    "InvalidMachineFiles",
				Field:   getNodeFieldYamlTag(node, idx, "MachineFiles"),
				Message: formatError(messages),
			})
		}
	}
	return result
}

func checkNodeExtensions(node Node, idx int, errs *Errors, warns *Warnings) (*Errors, *Warnings) {
	if len(node.Extensions) > 0 {
		warns.Append(&Warning{
			Kind:    "DeprecatedNodeExtensions",
			Field:   getNodeFieldYamlTag(node, idx, "Extensions"),
			Message: formatWarning("`extensions` is deprecated, please use `schematic.customization.systemExtensions` instead"),
		})
		var messages *multierror.Error
		extensions := map[string]struct{}{}

		for _, ext := range node.Extensions {
			if _, exists := extensions[ext.Image()]; exists {
				messages = multierror.Append(messages, fmt.Errorf("duplicate system extension %q", ext.Image()))
			}
			extensions[ext.Image()] = struct{}{}
		}

		if messages.ErrorOrNil() != nil {
			return errs.Append(&Error{
				Kind:    "InvalidNodeExtensions",
				Field:   getNodeFieldYamlTag(node, idx, "Extensions"),
				Message: formatError(messages),
			}), warns
		}
	}

	return errs, warns
}

func checkNodeSchematic(node Node, idx int, result *Errors) *Errors {
	var messages *multierror.Error
	extensions := map[string]struct{}{}
	if node.Schematic != nil {
		supportedExtensions := []string{
			"siderolabs/amd-ucode",
			"siderolabs/bnx2-bnx2x",
			"siderolabs/drbd",
			"siderolabs/gasket-driver",
			"siderolabs/gvisor",
			"siderolabs/hello-world-service",
			"siderolabs/i915-ucode",
			"siderolabs/intel-ucode",
			"siderolabs/iscsi-tools",
			"siderolabs/nut-client",
			"siderolabs/nvidia-container-toolkit",
			"siderolabs/nvidia-fabricmanager",
			"siderolabs/nvidia-open-gpu-kernel-modules",
			"siderolabs/qemu-guest-agent",
			"siderolabs/tailscale",
			"siderolabs/thunderbolt",
			"siderolabs/usb-modem-drivers",
			"siderolabs/zfs",
			"siderolabs/nonfree-kmod-nvidia",
		}
		for _, ext := range node.Schematic.Customization.SystemExtensions.OfficialExtensions {
			if !slices.Contains(supportedExtensions, ext) {
				messages = multierror.Append(messages, fmt.Errorf("%q is not a supported Talos extension", ext))
			}
			if _, exists := extensions[ext]; exists {
				messages = multierror.Append(messages, fmt.Errorf("duplicate system extension %q", ext))
			}
			extensions[ext] = struct{}{}
		}
	}

	if messages.ErrorOrNil() != nil {
		return result.Append(&Error{
			Kind:    "InvalidNodeSchematic",
			Field:   getNodeFieldYamlTag(node, idx, "Schematic"),
			Message: formatError(messages),
		})
	}

	return result
}

func checkNodeNameServers(node Node, idx int, result *Errors) *Errors {
	if len(node.Nameservers) > 0 {
		for _, ip := range node.Nameservers {
			if !validate.IsIP(ip) {
				e := fmt.Errorf("%q is not a valid list of IP addresses", node.Nameservers[:])
				return result.Append(&Error{
					Kind:    "InvalidNodeNameservers",
					Field:   getNodeFieldYamlTag(node, idx, "Nameservers"),
					Message: formatError(multierror.Append(e)),
				})
			}
		}
	}
	return result
}

func checkNodeNetworkInterfaces(node Node, idx int, result *Errors) *Errors {
	if len(node.NetworkInterfaces) > 0 {
		var (
			warnings []string
			messages *multierror.Error
		)
		bondedInterfaces := map[string]string{}
		bridgedInterfaces := map[string]string{}

		for _, device := range node.NetworkInterfaces {
			if device.Bond() != nil && device.Bridge() != nil {
				messages = multierror.Append(messages, fmt.Errorf("interface has both bridge and bond section set %q", device.Interface()))
			}

			if device.Bond() != nil {
				for _, iface := range device.Bond().Interfaces() {
					if otherIface, exists := bondedInterfaces[iface]; exists && otherIface != device.Interface() {
						messages = multierror.Append(messages, fmt.Errorf("interface %q is declared as part of two bonds: %q and %q", iface, otherIface, device.Interface()))
					}

					if bridgeIface, exists := bridgedInterfaces[iface]; exists {
						messages = multierror.Append(messages, fmt.Errorf("interface %q is declared as part of an interface and a bond: %q and %q", iface, bridgeIface, device.Interface()))
					}

					bondedInterfaces[iface] = device.Interface()
				}

				if len(device.Bond().Interfaces()) > 0 && len(device.Bond().Selectors()) > 0 {
					messages = multierror.Append(messages, fmt.Errorf("interface %q has both interfaces and selectors set", device.Interface()))
				}
			}

			if device.Bridge() != nil {
				for _, iface := range device.Bridge().Interfaces() {
					if otherIface, exists := bridgedInterfaces[iface]; exists && otherIface != device.Interface() {
						messages = multierror.Append(messages, fmt.Errorf("interface %q is declared as part of two bridges: %q and %q", iface, otherIface, device.Interface()))
					}

					if bondIface, exists := bondedInterfaces[iface]; exists {
						messages = multierror.Append(messages, fmt.Errorf("interface %q is declared as part of an interface and a bond: %q and %q", iface, bondIface, device.Interface()))
					}

					bridgedInterfaces[iface] = device.Interface()
				}
			}
			warn, err := v1alpha1.ValidateNetworkDevices(device, bondedInterfaces, v1alpha1.CheckDeviceInterface, v1alpha1.CheckDeviceAddressing, v1alpha1.CheckDeviceRoutes)
			warnings = append(warnings, warn...)
			messages = multierror.Append(messages, err)
			for _, w := range warnings {
				messages = multierror.Append(messages, fmt.Errorf(w))
			}
		}

		if messages.ErrorOrNil() != nil {
			return result.Append(&Error{
				Kind:    "InvalidNodeNetworkInterfaces",
				Field:   getNodeFieldYamlTag(node, idx, "NetworkInterfaces"),
				Message: formatError(messages),
			})
		}
	}

	return result
}

func checkNodeConfigPatches(node Node, idx int, result *Errors) *Errors {
	if len(node.ConfigPatches) > 0 {
		if !isRFC6902List(node.ConfigPatches) {
			e := fmt.Errorf("doesn't look like list of RFC6902 JSON patches")
			return result.Append(&Error{
				Kind:    "InvalidNodeConfigPatches",
				Field:   getNodeFieldYamlTag(node, idx, "ConfigPatches"),
				Message: formatError(multierror.Append(e)),
			})
		}
	}
	return result
}

func isDomain(domain string) bool {
	if domain == "" || len(domain)-strings.Count(domain, ".") > 255 {
		return false
	}
	return regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`).MatchString(domain)
}

func isDomainOrIP(domainIP string) bool {
	return isDomain(domainIP) || validate.IsIP(domainIP)
}

func isCIDRList(nets []string) bool {
	for _, net := range nets {
		if !validate.IsCIDR(net) {
			return false
		}
	}
	return true
}

func isSemVer(version string) bool {
	stripped := strings.TrimPrefix(version, "v")
	re := `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	if match, _ := regexp.MatchString(re, stripped); match {
		return true
	}
	return false
}

func isRFC6902List(input []map[string]interface{}) bool {
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

func formatError(e *multierror.Error) *multierror.Error {
	e.ErrorFormat = func(es []error) string {
		points := make([]string, len(es))
		for i, err := range es {
			points[i] = fmt.Sprintf("  * %s", err)
		}
		return strings.Join(points, "\n")
	}
	return e
}

func formatWarning(w string) string {
	return fmt.Sprintf("  * WARNING: %s", w)
}

func getNodeFieldYamlTag(node Node, idx int, fieldPath string) string {
	return "nodes[" + fmt.Sprintf("%v", idx) + "]." + getFieldYamlTag(node, fieldPath)
}

func getFieldYamlTag(v interface{}, fieldPath string) string {
	parts := strings.Split(fieldPath, ".")
	structValue := reflect.ValueOf(v)
	result := []string{}

	for i := 0; i < len(parts); i++ {
		fieldName := parts[i]
		field := structValue.FieldByName(fieldName)

		if !field.IsValid() {
			return fieldPath
		}

		yamlTag := ""
		if found, ok := structValue.Type().FieldByName(fieldName); ok {
			yamlTag = found.Tag.Get("yaml")
		}
		tagParts := strings.Split(yamlTag, ",")

		structValue = field
		result = append(result, tagParts[0])
	}
	return strings.Join(result, ".")
}
