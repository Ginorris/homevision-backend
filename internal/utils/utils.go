package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

// TODO we could add timestamps to the names to avoid conflicts
// SanitizeFilename removes or replaces characters that are unsafe for file names.
func SanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	name = reg.ReplaceAllString(name, "")
	name = strings.Trim(name, "_")
	return filepath.Clean(name)
}
