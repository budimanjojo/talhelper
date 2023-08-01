package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/budimanjojo/talhelper/pkg/patcher"
	"github.com/budimanjojo/talhelper/pkg/talos"
)

// GenerateConfig takes `TalhelperConfig` and path to encrypted `secretFile` and generates
// Talos `machineconfig` files and a `talosconfig` file in `outDir`.
// It returns an error, if any.
func GenerateConfig(c *config.TalhelperConfig, outDir, secretFile, mode string) error {
	var cfg []byte
	input, err := talos.NewClusterInput(c, secretFile)
	if err != nil {
		return err
	}

	for _, node := range c.Nodes {
		fileName := c.ClusterName + "-" + node.Hostname + ".yaml"
		cfgFile := outDir + "/" + fileName

		cfg, err = talos.GenerateNodeConfigBytes(&node, input)
		if err != nil {
			return err
		}

		if node.InlinePatch != nil {
			cfg, err = patcher.YAMLInlinePatcher(node.InlinePatch, cfg)
			if err != nil {
				return err
			}
		}

		if len(node.ConfigPatches) != 0 {
			cfg, err = patcher.YAMLPatcher(node.ConfigPatches, cfg)
			if err != nil {
				return err
			}
		}

		if len(node.Patches) != 0 {
			cfg, err = patcher.PatchesPatcher(node.Patches, cfg)
			if err != nil {
				return err
			}
		}

		if node.ControlPlane {
			cfg, err = patcher.YAMLInlinePatcher(c.ControlPlane.InlinePatch, cfg)
			if err != nil {
				return err
			}
			cfg, err = patcher.YAMLPatcher(c.ControlPlane.ConfigPatches, cfg)
			if err != nil {
				return err
			}
			cfg, err = patcher.PatchesPatcher(c.ControlPlane.Patches, cfg)
			if err != nil {
				return err
			}
		} else {
			cfg, err = patcher.YAMLInlinePatcher(c.Worker.InlinePatch, cfg)
			if err != nil {
				return err
			}
			cfg, err = patcher.YAMLPatcher(c.Worker.ConfigPatches, cfg)
			if err != nil {
				return err
			}
			cfg, err = patcher.PatchesPatcher(c.Worker.Patches, cfg)
			if err != nil {
				return err
			}
		}

		err = talos.ValidateConfigFromBytes(cfg, mode)
		if err != nil {
			return err
		}

		cfg, err = talos.ReEncodeTalosConfig(cfg)
		if err != nil {
			return err
		}

		err = dumpFile(cfgFile, cfg)
		if err != nil {
			return err
		}

		fmt.Printf("generated config for %s in %s\n", node.Hostname, cfgFile)
	}

	clientCfg, err := talos.GenerateClientConfigBytes(c, input)
	if err != nil {
		return err
	}

	fileName := "talosconfig"
	err = dumpFile(outDir + "/" + fileName, clientCfg)
	if err != nil {
		return err
	}

	fmt.Printf("generated client config in %s\n", outDir + "/" + fileName)

	return nil
}

// dumpFile creates file in `path` and dumps the content of bytes into
// the path. It returns an error, if any.
func dumpFile(path string, file []byte) error {
	dirName := filepath.Dir(path)

	_, err := os.Stat(dirName)
	if err != nil {
		err := os.MkdirAll(dirName, 0700)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path, file, 0600)
	if err != nil {
		return err
	}

	return nil
}
