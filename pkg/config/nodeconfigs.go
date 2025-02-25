package config

import (
	"reflect"
)

func (node *Node) OverrideGlobalCfg(cfg NodeConfigs) *Node {
	node.NodeConfigs = mergeNodeConfigs(node.NodeConfigs, cfg, node.OverridePatches, node.OverrideExtraManifests, node.OverrideExtraCertSANs)

	return node
}

func mergeNodeConfigs(patch, src NodeConfigs, overridePatches, overrideExtraManifest, overrideExtraCertSANs bool) NodeConfigs {
	if len(src.Patches) > 0 && !overridePatches {
		// global patches should get applied first
		// https://github.com/budimanjojo/talhelper/issues/388
		patch.Patches = append(src.Patches, patch.Patches...)
	}
	if len(src.ExtraManifests) > 0 && !overrideExtraManifest {
		patch.ExtraManifests = append(patch.ExtraManifests, src.ExtraManifests...)
	}
	if len(src.ExtraCertSANs) > 0 && !overrideExtraCertSANs {
		patch.ExtraCertSANs = append(patch.ExtraCertSANs, src.ExtraCertSANs...)
	}

	patchValue := reflect.ValueOf(patch)
	srcValue := reflect.ValueOf(src)

	result := reflect.New(patchValue.Type()).Elem()

	for i := range patchValue.NumField() {
		patchField := patchValue.Field(i)
		srcField := srcValue.Field(i)

		if !patchField.IsZero() {
			result.Field(i).Set(patchField)
		} else {
			result.Field(i).Set(srcField)
		}
	}

	return result.Interface().(NodeConfigs)
}
