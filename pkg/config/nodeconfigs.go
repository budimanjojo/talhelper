package config

import (
	"reflect"
)

func (node *Node) OverrideGlobalCfg(cfg NodeConfigs) *Node {
	node.NodeConfigs = mergeNodeConfigs(node.NodeConfigs, cfg, node.OverridePatches, node.OverrideExtraManifests)

	return node
}

func mergeNodeConfigs(patch, src NodeConfigs, overridePatches, overrideExtraManifest bool) NodeConfigs {
	if len(src.Patches) > 0 && !overridePatches {
		patch.Patches = append(patch.Patches, src.Patches...)
	}
	if len(src.ExtraManifests) > 0 && !overrideExtraManifest {
		patch.ExtraManifests = append(patch.ExtraManifests, src.ExtraManifests...)
	}

	patchValue := reflect.ValueOf(patch)
	srcValue := reflect.ValueOf(src)

	result := reflect.New(patchValue.Type()).Elem()

	for i := 0; i < patchValue.NumField(); i++ {
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
