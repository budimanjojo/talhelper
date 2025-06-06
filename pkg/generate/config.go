package generate

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/budimanjojo/talhelper/v3/pkg/patcher"
	"github.com/budimanjojo/talhelper/v3/pkg/talos"
	"github.com/budimanjojo/talhelper/v3/pkg/templating"
	taloscfg "github.com/siderolabs/talos/pkg/machinery/config"
)

// GenerateConfig takes `TalhelperConfig` and path to encrypted `secretFile` and generates
// Talos `machineconfig` files and a `talosconfig` file in `outDir`.
// It returns an error, if any.
func GenerateConfig(c *config.TalhelperConfig, dryRun bool, outDir, secretFile, mode string, offlineMode bool, disableNodesSection bool) error {
	input, err := talos.NewClusterInput(c, secretFile, mode)
	if err != nil {
		return err
	}

	for _, node := range c.Nodes {
		var rawcfg taloscfg.Provider

		fileName, err := node.GetOutputFileName(c)
		if err != nil {
			return err
		}
		cfgFile := outDir + "/" + fileName
		slog.Debug(fmt.Sprintf("generating %s for node %s", cfgFile, node.Hostname))

		rawcfg, err = talos.GenerateNodeConfig(&node, input, c.GetImageFactory(), offlineMode)
		if err != nil {
			return err
		}

		// this is needed because the upstream cluster input doesn't handle inline manifests and some others so we need to do it ourselves
		if len(c.ClusterInlineManifests) > 0 {
			rawcfg.RawV1Alpha1().ClusterConfig.ClusterInlineManifests = *c.ClusterInlineManifests.GetIMs()
		}

		cfg, err := rawcfg.Bytes()
		if err != nil {
			return err
		}

		if len(node.Patches) != 0 {
			slog.Debug(fmt.Sprintf("applying node specific patches to %s", node.Hostname))
			cfg, err = patcher.PatchesPatcher(node.Patches, cfg)
			if err != nil {
				return err
			}
		}

		if len(c.Patches) > 0 {
			slog.Debug(fmt.Sprintf("applying global patches to %s", node.Hostname))
			cfg, err = patcher.PatchesPatcher(c.Patches, cfg)
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

		if node.IngressFirewall != nil {
			slog.Debug(fmt.Sprintf("generating machine firewall config for %s", node.Hostname))
			nc, err := talos.GenerateNetworkConfigBytes(node.IngressFirewall)
			if err != nil {
				return err
			}
			cfg = append(cfg, nc...)
		}

		if len(node.ExtensionServices) > 0 {
			slog.Debug(fmt.Sprintf("generating machine extension service config for %s", node.Hostname))
			ext, err := talos.GenerateExtensionServicesConfigBytes(node.ExtensionServices)
			if err != nil {
				return err
			}
			cfg = append(cfg, ext...)
		}

		if len(node.Volumes) > 0 {
			slog.Debug(fmt.Sprintf("generating volume config for %s", node.Hostname))
			vc, err := talos.GenerateVolumeConfigBytes(node.Volumes, mode)
			if err != nil {
				return err
			}
			cfg = append(cfg, vc...)
		}

		if len(node.UserVolumes) > 0 {
			slog.Debug(fmt.Sprintf("generating user volume config for %s", node.Hostname))
			uvc, err := talos.GenerateUserVolumeConfigBytes(node.UserVolumes, mode)
			if err != nil {
				return err
			}
			cfg = append(cfg, uvc...)
		}

		if len(node.ExtraManifests) > 0 {
			slog.Debug(fmt.Sprintf("generating extra manifests for %s", node.Hostname))
			content, err := combineExtraManifests(node.ExtraManifests, cfg)
			if err != nil {
				return err
			}
			cfg = append(cfg, content...)
		}

		if !dryRun {
			slog.Debug(fmt.Sprintf("dumping machineconfig file for %s to %s", node.Hostname, cfgFile))
			err = dumpFile(cfgFile, cfg)
			if err != nil {
				return err
			}

			fmt.Printf("generated config for %s in %s\n", node.Hostname, cfgFile)
		} else {
			slog.Debug("showing diff from previous run")
			absCfgFile, err := filepath.Abs(cfgFile)
			if err != nil {
				return err
			}

			before, err := getFileContent(absCfgFile)
			if err != nil {
				return err
			}

			diff := computeDiff(absCfgFile, before, string(cfg))
			if diff != "" {
				fmt.Println(diff)
			} else {
				fmt.Printf("no changes found on %s\n", cfgFile)
			}
		}
	}

	if !dryRun {
		clientCfg, err := talos.GenerateClientConfigBytes(c, input, disableNodesSection)
		if err != nil {
			return err
		}

		fileName := "talosconfig"

		slog.Debug(fmt.Sprintf("dumping talosconfig file to %s", outDir+"/"+fileName))
		err = dumpFile(outDir+"/"+fileName, clientCfg)
		if err != nil {
			return err
		}

		fmt.Printf("generated client config in %s\n", outDir+"/"+fileName)
	}

	return nil
}

// getFileContent returns content of file as string.
// It also returns an error, if any
func getFileContent(path string) (string, error) {
	content, err := getFileContentByte(path)
	return string(content), err
}

// getFileContentByte returns content of file. It also returns an error,
// if any
func getFileContentByte(path string) ([]byte, error) {
	if _, osErr := os.Stat(path); osErr == nil {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return content, nil
	} else if errors.Is(osErr, os.ErrNotExist) {
		return nil, nil
	} else {
		return nil, osErr
	}
}

// combineExtraManifests takes list of filepaths, parse go template using data from
// main config manifest, combines them into a single file in bytes with `---\n` prepended.
// It also returns an error, if any
func combineExtraManifests(extraFiles []string, templateData []byte) ([]byte, error) {
	var result [][]byte
	dataCfg, err := talos.LoadTalosConfig(templateData)
	if err != nil {
		return nil, err
	}

	for _, file := range extraFiles {
		content, err := getFileContentByte(strings.TrimPrefix(file, "@"))
		if err != nil {
			return nil, err
		}

		content, err = templating.RenderTemplate[[]byte](string(content), dataCfg.RawV1Alpha1())
		if err != nil {
			return nil, err
		}
		result = append(result, content)
	}
	return talos.CombineYamlBytes(result), nil
}

// computeDiff returns diff between before and after string
// using Myers diff algorithm
func computeDiff(path, before, after string) string {
	edits := myers.ComputeEdits(span.URIFromPath(path), before, after)
	diff := gotextdiff.ToUnified("a"+path, "b"+path, before, edits)
	return fmt.Sprint(diff)
}

// dumpFile creates file in `path` and dumps the content of bytes into
// the path. It returns an error, if any.
func dumpFile(path string, file []byte) error {
	dirName := filepath.Dir(path)

	_, err := os.Stat(dirName)
	if err != nil {
		err := os.MkdirAll(dirName, 0o700)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path, file, 0o600)
	if err != nil {
		return err
	}

	return nil
}
