// cache.go
package main

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// purgeCache purges the cache file.
func purgeCache() {
	// Remove the cache file
	path := getCachePath()
	os.Remove(path)
	// Write blank cache file
	err := writeCache(&TalosVersionTags{})
	if err != nil {
		log.Errorf("error writing empty cache file: %s", err)
	}
}

// getCachePath returns the path to the cache directory based on the OS.
func getCachePath() string {
	// Get the current user
	currentUser, _ := user.Current()
	// Get the home directory of the user
	homeDir := currentUser.HomeDir
	// Define the cache directory based on the operating system
	var cacheDir string
	switch os := runtime.GOOS; os {
	case "windows":
		cacheDir = filepath.Join(homeDir, "AppData", "Local", "Cache", "talhelper")
	case "darwin": // macOS
		cacheDir = filepath.Join(homeDir, "Library", "Caches", "talhelper")
	default: // Linux and other Unix-like systems
		cacheDir = getEnv("XDG_CACHE_HOME", filepath.Join(homeDir, ".cache", "talhelper"))
	}

	return cacheDir
}

// checkCache checks if the cache file exists.
func checkCache() bool {
	// If the cache file doesn't exist
	if _, err := os.Stat(TSEHelperTalosExtensionsCachePath); os.IsNotExist(err) {
		return false
	}

	// else, it does exist
	return true
}

// loadCache loads the cache file into a TalosVersionTags struct.
func loadCache(versionTags *TalosVersionTags) (TalosVersionTags, error) {
	var talosVersionTags TalosVersionTags

	// Read the cache file
	log.Debugf("reading from cache file: %s", TSEHelperTalosExtensionsCachePath)
	if tags, err := os.ReadFile(TSEHelperTalosExtensionsCachePath); err != nil {
		return talosVersionTags, err
	} else if err := json.Unmarshal(tags, &talosVersionTags); err != nil {
		return talosVersionTags, err
	}
	log.Tracef("cache file contents: %s", talosVersionTags)
	return talosVersionTags, nil
}

// writeCache writes the TalosVersionTags struct to the cache file.
func writeCache(versionTags *TalosVersionTags) error {
	// Check if the cache directory exists or create
	if _, err := os.Stat(TSEHelperTalosExtensionsCachePath); os.IsNotExist(err) {
		if err := os.MkdirAll(TSEHelperTalosExtensionsCacheDir, 0755); err != nil {
			return err
		}
	}

	// Marshal the data with indentation
	bytes, err := json.MarshalIndent(versionTags, "", "    ")
	if err != nil {
		return err
	}

	// Write the pretty-printed JSON to the cache file
	if err := os.WriteFile(TSEHelperTalosExtensionsCachePath, bytes, 0644); err != nil {
		return err
	}

	return nil
}
