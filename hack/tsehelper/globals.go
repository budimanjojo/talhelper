package main

import (
	"os"
	"path/filepath"
)

// Constants for default values
const (
	DefaultStringRegex         = `^(?P<registry>[\w\.\-0-9]+)\/(?P<org>[\w\.\-0-9]+)\/(?P<repo>[\w\.\-0-9]+):(?P<tag>[\w\.\-0-9]+)(@sha256:(?P<shasum>[a-f0-9]+))?$`
	DefaultTalosExtensionsRepo = "ghcr.io/siderolabs/extensions"
	DefaultOutputFilename      = "talos-extensions.json"
	DefaultTalosOverlaysRepo   = "ghcr.io/siderolabs/overlays"
)

// Global variables
var (
	TSEHelperTalosExtensionsRegex      = getEnv("TSEHELPER_REGEX_OVERRIDE", DefaultStringRegex)
	TSEHelperTalosExtensionsRepository = getEnv("TSEHELPER_TALOS_EXTENSIONS_REPO", DefaultTalosExtensionsRepo)
	TSEHelperTalosExtensionsCacheDir   = getEnv("TSEHELPER_TALOS_EXTENSIONS_CACHE_DIR", getCachePath())
	TSEHelperTalosExtensionsCacheFile  = getEnv("TSEHELPER_TALOS_EXTENSIONS_CACHE_FILE", DefaultOutputFilename)
	TSEHelperTalosExtensionsCachePath  = filepath.Join(TSEHelperTalosExtensionsCacheDir, TSEHelperTalosExtensionsCacheFile)
	TSEHelperTalosOverlaysRepository   = getEnv("TSEHELPER_TALOS_OVERLAYS_REPO", DefaultTalosOverlaysRepo)
)

// getEnv is a utility function to get an environment variable if it exists, otherwise it will return the fallback value.
func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}
