// system_extensions.go
package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	log "github.com/sirupsen/logrus"
)

// Functions related to fetching system extensions
// getSystemExtensions fetches the system extensions for each Talos version and updates the TalosVersionTags struct.
// This function uses goroutines and channels to limit concurrency and collect errors.
func getSystemExtensions(tags *TalosVersionTags) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	maxWorkers := runtime.GOMAXPROCS(0)
	log.Debugf("maxWorkers: %d", maxWorkers)
	semaphore := make(chan struct{}, maxWorkers)
	errors := make(chan error, len(tags.Versions))

	log.Tracef("tags passed to goroutines: %v", tags)

	// Loop through the tags
	for i := range tags.Versions {
		wg.Add(1) // Increment the WaitGroup counter
		log.Tracef("working on: %s", tags.Versions[i])

		go func(i int) {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes

			// Acquire a slot from the semaphore (limiting concurrency)
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Release the slot when the goroutine completes

			sysExt := &tags.Versions[i]
			imageName := TSEHelperTalosExtensionsRepository + ":" + sysExt.Version
			desc, err := crane.Get(imageName)

			if err != nil {
				errors <- fmt.Errorf("error getting image %s: %s", imageName, err)
				return
			}

			var img v1.Image
			if desc.MediaType.IsSchema1() {
				img, err = desc.Schema1()
				if err != nil {
					errors <- fmt.Errorf("error getting schema1 for image %s: %s", imageName, err)
					return
				}
			} else {
				img, err = desc.Image()
				if err != nil {
					errors <- fmt.Errorf("error getting image for %s: %s", imageName, err)
					return
				}
			}
			var tarBuffer bytes.Buffer
			err = crane.Export(img, &tarBuffer)
			if err != nil {
				errors <- fmt.Errorf("error exporting image %s: %s", imageName, err)
				return
			}
			extensions, err := processTarArchive(tarBuffer.Bytes())
			if err != nil {
				errors <- fmt.Errorf("error processing tar archive: %s", err)
				return
			}
			// Remove empty strings
			var nonEmptyExtensions []string
			for _, ext := range extensions {
				if ext != "" {
					nonEmptyExtensions = append(nonEmptyExtensions, ext)
				}
			}
			// Update sysExt within the critical section to avoid race conditions
			mu.Lock()
			sysExt.SystemExtensions = nonEmptyExtensions
			mu.Unlock()
		}(i)
	}

	// Close the errors channel when all goroutines have completed
	go func() {
		wg.Wait()
		close(errors)
	}()

	// Collect errors into one aggregated error string and return it.
	errString := []string{}
	for err := range errors {
		errString = append(errString, err.Error())
	}
	if len(errString) > 0 {
		return fmt.Errorf(strings.Join(errString, "\n"))
	}

	// Wait for all goroutines to complete
	wg.Wait()
	log.Trace("finished goroutines")

	return nil
}
