package watch

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Simple helper package to watch files for changes.

// Files will watch files matching filePatterns under the directory
// dirPath. Returns only if/when a file change is detected.
func Files(dirPath string, filePatterns []string) {

	infos := findFiles(dirPath, filePatterns)

	if len(infos) == 0 {
		log.Fatalf("failed to find any files to watch")
	}

	log.Printf("watching %d files matching %v", len(infos), filePatterns)

	for {
		time.Sleep(500 * time.Millisecond)
		if checkFiles(infos) {
			return
		}
	}
}

// checkFiles returns true if a file has changed, false if nothing has changed.
// It goes thru the list of files being watched, checking each to see if the
// modified timestamp has changed.
func checkFiles(infoMap map[string]os.FileInfo) bool {

	for path, info := range infoMap {

		newInfo, err := os.Stat(path)
		if err != nil {
			log.Printf("%s was removed", path)
			return true
		}

		if newInfo.ModTime().After(info.ModTime()) {
			log.Printf("%s was modified", path)
			return true
		}
	}

	return false
}

// findFiles takes a dir path and a slice of patterns. it finds all files under
// dir matching patterns, and returns info on all the found files.
func findFiles(dir string, patterns []string) map[string]os.FileInfo {

	found := map[string]os.FileInfo{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() != "." && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			return nil
		}

		for _, p := range patterns {
			match, err := filepath.Match(p, info.Name())
			if err != nil {
				return err
			}
			if match {
				found[path] = info
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("failed to find files in dir %s: %v", dir, err)
	}

	return found
}
