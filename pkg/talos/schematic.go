package talos

import (
	"strings"

	"github.com/siderolabs/image-factory/pkg/schematic"
)

func GetInstallerURL(cfg *schematic.Schematic, registryURL, version string) (string, error) {
	id, err := cfg.ID()
	if err != nil {
		return "", err
	}
	return ensureSlashSuffix(registryURL) + id + ":" + version, nil
}

func GetISOURL(cfg *schematic.Schematic, registryURL, version, mode, arch string) (string, error) {
	id, err := cfg.ID()
	if err != nil {
		return "", err
	}
	return ensureSlashSuffix(registryURL) + ensureSlashSuffix(id) + ensureSlashSuffix(version) + mode + "-" + arch + ".iso", nil
}

func ensureSlashSuffix(s string) string {
	if strings.HasSuffix(s, "/") {
		return s
	}
	return s + "/"
}
