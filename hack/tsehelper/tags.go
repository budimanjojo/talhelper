package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"tsehelper/pkg/versiontags"

	"github.com/google/go-containerregistry/pkg/crane"
	log "github.com/sirupsen/logrus"
)

// getMissingTags fetches all tags for a given repository and return a TalosVersionsTags struct of all tags not already in the cache or an error.
func getMissingTags(cachedTags versiontags.TalosVersionTags) (versiontags.TalosVersionTags, error) {
	tagsToAppend := versiontags.TalosVersionTags{}

	// Fetch the tags
	log.Debugf("calling registry docker://%s...", TSEHelperTalosExtensionsRepository)
	upstreamTags, err := crane.ListTags(TSEHelperTalosExtensionsRepository)
	if err != nil {
		return cachedTags, err
	}

	// Loop through the tags
	for _, tag := range upstreamTags {
		// Skip anything that doesn't start with v
		if !strings.HasPrefix(tag, "v") {
			log.Tracef("skipping tag %s", tag)
			continue
		}
		// Skip any tag that's already present
		if cachedTags.Contains(tag) {
			continue
		}

		// Add any new tags to the list
		log.Debugf("adding new tag %s", tag)
		tagsToAppend.Versions = append(tagsToAppend.Versions, versiontags.TalosVersion{Version: tag})
	}

	// Sort the list
	log.Debugf("finalizing list of tags to append: %s", tagsToAppend)
	sort.Sort(tagsToAppend)

	return tagsToAppend, nil
}

// getMissingVersions checks if a newer version is available and returns a TalosVersionTags struct of all missing versions.
func getMissingVersions(versionsTags *versiontags.TalosVersionTags) versiontags.TalosVersionTags {
	// Check if the cache file exists
	if !checkCache() {
		// Load the cache file
		_, err := loadCache(versionsTags)
		if err != nil {
			log.Errorf("error loading cache: %s", err)
			os.Exit(1)
		}
	}

	// Fetch the missing tags
	tags, err := getMissingTags(*versionsTags)
	if err != nil {
		fmt.Printf("Error fetching tags: %s\n", err)
		os.Exit(1)
	}

	return tags
}

// cleanString takes a string and returns a cleaned string based on the flags passed.
func cleanString(line string) string {
	// Create a new regexp from the TalHelperTalosExtensionsRegex
	regexp := regexp.MustCompile(TSEHelperTalosExtensionsRegex)

	// Find the sub-matches
	matches := regexp.FindStringSubmatch(line)
	if len(matches) == 0 {
		log.Tracef("no matches found for line: %s", line)
		return ""
	}
	log.Tracef("regexp matches: %s", matches)

	// Map results to capture group names
	result := make(map[string]string)
	for i, name := range regexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	// If all flags are set or minimal is set, return the minimal json
	if (trimRegistry && trimSha256 && trimTag) || minimal {
		log.Tracef("returning minimal json")
		return fmt.Sprintf(result["org"] + "/" + result["repo"])
	}

	if trimRegistry && trimSha256 {
		log.Tracef("returning trimmed registry and sha256")
		return fmt.Sprintf(result["org"] + "/" + result["repo"] + ":" + result["tag"])
	} else if trimRegistry && !trimSha256 {
		if trimTag {
			log.Tracef("returning trimmed registry and tag")
			return fmt.Sprintf(result["org"] + "/" + result["repo"] + "@sha256:" + result["shasum"])
		}
		log.Tracef("returning trimmed registry")
		return fmt.Sprintf(result["org"] + "/" + result["repo"] + ":" + result["tag"] + "@sha256:" + result["shasum"])
	} else if !trimRegistry && trimSha256 {
		if trimTag {
			log.Tracef("returning trimmed sha256 and tag")
			return fmt.Sprintf(result["registry"] + "/" + result["org"] + "/" + result["repo"])
		}
		log.Tracef("returning trimmed sha256")
		return fmt.Sprintf(result["registry"] + "/" + result["org"] + "/" + result["repo"] + ":" + result["tag"])
	} else {
		if trimTag {
			log.Tracef("returning trimmed tag")
			return fmt.Sprintf(result["registry"] + "/" + result["org"] + "/" + result["repo"] + "@sha256:" + result["shasum"])
		}
		log.Tracef("returning full string")
		return fmt.Sprintf(result["registry"] + "/" + result["org"] + "/" + result["repo"] + ":" + result["tag"] + "@sha256:" + result["shasum"])
	}
}
