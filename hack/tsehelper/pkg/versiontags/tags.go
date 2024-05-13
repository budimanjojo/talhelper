package versiontags

import (
	"golang.org/x/mod/semver"
)

// TalosVersion is a struct that holds the Talos version and list of available Talos System Extensions.
type TalosVersion struct {
	Version          string    `json:"version"`
	SystemExtensions []string  `json:"systemExtensions"`
	Overlays         []Overlay `json:"overlays" yaml:"overlays"`
}

// Overlay is a struct that holds the Talos Overlay information.
type Overlay struct {
	Image  string `json:"image" yaml:"image"`
	Name   string `json:"name" yaml:"name"`
	Digest string `json:"digest" yaml:"digest"`
}

// TalosVersionTags is a struct that holds the list of TalosVersionTags for each Talos version returned by the registry.
type TalosVersionTags struct {
	Versions []TalosVersion `json:"versions"`
}

// Implement Contains on TalosVersionsTags.Versions
func (v TalosVersionTags) Contains(s string) bool {
	for _, a := range v.Versions {
		if a.Version == s {
			return true
		}
	}

	return false
}

// Implement Len interface on TalosVersionsTags.Versions
func (v TalosVersionTags) Len() int {
	return len(v.Versions)
}

func (v TalosVersionTags) Less(i, j int) bool {
	return semver.Compare(v.Versions[i].Version, v.Versions[j].Version) < 0
}

func (v TalosVersionTags) Swap(i, j int) {
	v.Versions[i], v.Versions[j] = v.Versions[j], v.Versions[i]
}

func (v TalosVersionTags) SliceIndex(s string) int {
	for i, a := range v.Versions {
		if a.Version == s {
			return i
		}
	}
	return -1
}
