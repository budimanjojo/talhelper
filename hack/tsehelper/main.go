package main

import (
	"os"
	"sort"

	"tsehelper/pkg/versiontags"

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
		log.Debug("cache purged successfully")
	}

	var tags versiontags.TalosVersionTags

	// Check if the cache file exists
	if checkCache() {
		// Load the cache file
		log.Debugf("cache exists, loading...")
		tags, err := loadCache(&tags)
		if err != nil {
			log.Errorf("error loading cache: %s", err)
			os.Exit(1)
		}
		if skipUpdate {
			log.Debug("skipping upstream check")
		} else {
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
				err := writeCache(&tags)
				if err != nil {
					log.Errorf("error writing cache file: %s", err)
					os.Exit(1)
				}
				log.Debug("missing versions added and cache written successfully")
			}
		}
	} else {
		if skipUpdate {
			log.Debug("skipping upstream check")
		} else {
			// No cache, fetch all tags
			log.Debugf("cache not found, fetching all tags...")
			// this is needed so `go` doesn't think `tags` is a
			// new variable inside the scope of this block
			var err error
			tags, err = getMissingTags(tags)
			if err != nil {
				log.Errorf("error fetching tags: %s", err)
				os.Exit(1)
			}

			if len(tags.Versions) == 0 {
				log.Errorf("no versions found")
				os.Exit(1)
			}

			// Parse the tags for system extensions
			if err := getSystemExtensions(&tags); err != nil {
				log.Errorf("error parsing tags for system extensions: %s", err)
				os.Exit(1)
			}
		}

		// Sort the cache
		sort.Sort(&tags)
		// Save the cache file
		err := writeCache(&tags)
		if err != nil {
			log.Errorf("error writing cache file: %s", err)
			os.Exit(1)
		}

		log.Debug("missing tags fetched, extensions parsed, and cache written successfully")
	}

	log.Debugf("finished updating and loading cache")

	log.Debug("preparing to output to stdout...")
	// Read the cache file
	tags, err := loadCache(&tags)
	if err != nil {
		log.Errorf("error loading cache: %s", err)
		os.Exit(1)
	}

	// Generate the output
	if len(tags.Versions) == 0 {
		log.Info("no versions found")
		os.Exit(1)
	}

	if err := os.WriteFile(output, generateOutput(tags), 0o755); err != nil {
		log.Errorf("failed to write file to %s: %s", output, err)
		os.Exit(1)
	}

	// Exit successfully!
	os.Exit(0)
}
