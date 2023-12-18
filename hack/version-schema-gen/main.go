package main

import (
	"encoding/json"
	"fmt"

	"os"
	"sort"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Initialize logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	if level, err := log.ParseLevel(getEnv("LOG_LEVEL", "info")); err == nil {
		log.SetLevel(level)
		log.Debugf("LOG_LEVEL: %s", level)
	}
}

// main does the business
func main() {
	// Check if purge flag is set
	log.Debugf("purge cache: %t", purge)
	if purge {
		purgeCache()
		log.Info("cache purged successfully")
	}

	var tags TalosVersionTags

	// Check if the cache file exists
	if checkCache() {
		// Load the cache file
		log.Debugf("cache exists, loading...")
		tags, err := loadCache(&tags)
		if err != nil {
			log.Errorf("error loading cache: %s", err)
			os.Exit(1)
		}

		// Check for missing talos versions
		if missingVersions := getMissingVersions(&tags); len(missingVersions.Versions) > 0 {
			// Populate missing versions with system extensions
			if err := getSystemExtensions(&missingVersions); err != nil {
				log.Errorf("error parsing tags for system extensions: %s", err)
				os.Exit(1)
			}

			// Append the newly resolved versions to the cache
			tags.Versions = append(tags.Versions, missingVersions.Versions...)

			// Sort the cache
			sort.Sort(&tags)
			// Write the cache file
			writeCache(&tags)
			log.Info("missing versions added and cache written successfully")
		}
	} else {
		// No cache, fetch all tags
		log.Debugf("cache not found, fetching all tags...")
		tags, err := getMissingTags(tags)
		if err != nil {
			log.Errorf("error fetching tags: %s", err)
			os.Exit(1)
		}

		// Parse the tags for system extensions
		if err := getSystemExtensions(&tags); err != nil {
			log.Errorf("error parsing tags for system extensions: %s", err)
			os.Exit(1)
		}

		// Sort the cache
		sort.Sort(&tags)
		// Save the cache file
		writeCache(&tags)
		log.Info("tags fetched, system extensions parsed, and cache written successfully")
	}
	log.Debugf("finished updating and loading cache")

	log.Debug("preparing to output to stdout...")
	// Read the cache file
	tags, err := loadCache(&tags)
	if err != nil {
		log.Errorf("error loading cache: %s", err)
		os.Exit(1)
	}

	log.Debug("cleaning strings...")
	// For each Talos version, call cleanString for each SystemExtensions
	for i, version := range tags.Versions {
		for j, extension := range version.SystemExtensions {
			tags.Versions[i].SystemExtensions[j] = cleanString(extension)
		}
	}

	log.Debug("preparing to marshal to JSON...")
	// Marshal tags to JSON with indentation
	bytes, err := json.MarshalIndent(tags, "", "    ")
	if err != nil {
		log.Errorf("error marshalling tags: %s", err)
		os.Exit(1)
	}
	log.Debug("printing JSON to stdout...")
	fmt.Println(string(bytes))

	// Exit successfully!
	os.Exit(0)
}
