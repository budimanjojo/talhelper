package main

import (
	"encoding/json"
	"os"

	"tsehelper/pkg/versiontags"

	log "github.com/sirupsen/logrus"
)

type Versions struct {
	Versions []string `json:"versions"`
}

func generateOutput(givenVersionTags versiontags.TalosVersionTags) []byte {
	// Log all flags in trace.
	log.Tracef("minimal: %t", minimal)

	log.Tracef("onlyVersions: %t", onlyVersions)
	log.Tracef("specific Version: %s", specificVersion)

	log.Tracef("trimRegistry: %t", trimRegistry)
	log.Tracef("trimSha256: %t", trimSha256)
	log.Tracef("trimTag: %t", trimTag)

	log.Debug("cleaning strings...")
	// For each Talos version, call cleanString for each SystemExtensions
	if onlyVersions {
		var vers Versions
		for _, val := range givenVersionTags.Versions {
			vers.Versions = append(vers.Versions, val.Version)
		}
		// Marshal tags to JSON with indentation
		bytes, err := json.MarshalIndent(vers, "", "    ")
		if err != nil {
			log.Errorf("error marshalling tags: %s", err)
			os.Exit(1)
		}
		return bytes
	} else if specificVersion != "" {
		var theOne versiontags.TalosVersionTags
		for _, val := range givenVersionTags.Versions {
			if val.Version == specificVersion {
				for i, extension := range val.SystemExtensions {
					val.SystemExtensions[i] = cleanString(extension)
				}
				theOne.Versions = append(theOne.Versions, val)
			}
		}
		if len(theOne.Versions) == 0 {
			log.Errorf("no version found for specified version %s", specificVersion)
			os.Exit(1)
		}
		// Marshal tags to JSON with indentation
		bytes, err := json.MarshalIndent(theOne, "", "    ")
		if err != nil {
			log.Errorf("error marshalling tags: %s", err)
			os.Exit(1)
		}
		return bytes
	}

	for i, val := range givenVersionTags.Versions {
		for j, extension := range val.SystemExtensions {
			givenVersionTags.Versions[i].SystemExtensions[j] = cleanString(extension)
		}
	}

	log.Debug("preparing to marshal to JSON...")
	// Marshal tags to JSON with indentation
	bytes, err := json.MarshalIndent(givenVersionTags, "", "    ")
	if err != nil {
		log.Errorf("error marshalling tags: %s", err)
		os.Exit(1)
	}
	log.Debugf("writing JSON to %s...", output)
	return bytes
}
