package util

import (
	"os"
	"path/filepath"
)

func WhereAmI() string {
	cwd, _ := os.Executable()
	if cwd == "" {
		cwd = "go_2_hell:*"
	}
	return cwd
}

func CurrentRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree looking for a file that indicates the project root
	for {
		configFilePath := filepath.Join(currentDir, "config.yaml")
		if _, err := os.Stat(configFilePath); err == nil {
			return currentDir, nil
		}

		// Move up one directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break // Reached the root without finding the file
		}
		currentDir = parentDir
	}

	return "", nil
}
