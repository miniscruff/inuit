package internal

import (
	"fmt"
	"os"
	"strings"
)

func ExistingScenes() ([]string, error) {
	var dirs []string

	entries, err := os.ReadDir(InternalDir)
	if err != nil {
		return dirs, fmt.Errorf("failure to read scenes: %w", err)
	}

	for _, dir := range entries {
		n := dir.Name()
		if n == AssetsFile || n == ContentsFile || n == MetadataFile {
			continue
		}

		dirs = append(dirs, strings.TrimSuffix(n, ".json"))
	}

	return dirs, nil
}
