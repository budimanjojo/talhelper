package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

var (
	factoryUrl        = "https://factory.talos.dev"
	defaultExtensions = []string{
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
)

type Version struct {
	Name string
}

func getVersions() []string {
	url := factoryUrl + "/versions"
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error doing GET request to %s: %s", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatalf("%s returned %s", url, response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body from %s: %s", url, err)
	}

	var versions []string

	if err := json.Unmarshal(body, &versions); err != nil {
		log.Fatalf("Error unmarshalling JSON: %s", err)
	}

	sort.Strings(versions)

	return versions
}

func getExtensions(versions []string) map[string][]string {
	result := make(map[string][]string)
	result["default"] = defaultExtensions
	for _, version := range versions {
		url := factoryUrl + "/version/" + version + "/extensions/official"
		response, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error doing GET request to %s, %s", url, err)
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusInternalServerError {
			result[version] = defaultExtensions
			continue
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Error reading response body from %s, %s", url, err)
		}

		var r []Version

		if err := json.Unmarshal(body, &r); err != nil {
			log.Fatalf("Error unmarshalling JSON: %s %s", body, err)
		}

		for _, a := range r {
			result[version] = append(result[version], a.Name)
		}

	}
	return result
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalf("no output file")
	}

	if _, err := os.Stat(filepath.Dir(flag.Arg(0))); os.IsNotExist(err) {
		log.Fatalf("%s", err)
	}

	versions := getVersions()
	result := getExtensions(versions)

	final, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := os.WriteFile(flag.Arg(0), final, 0o755); err != nil {
		log.Fatal(err)
	}
}
