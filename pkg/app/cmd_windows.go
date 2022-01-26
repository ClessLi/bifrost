package app

import (
	"strings"
)

// FormatBasename is formatted as an executable file name under different
// operating systems according to the given name.
func FormatBasename(basename string) string {
	// Make case-insensitive and strip executable suffix if present
	basename = strings.ToLower(basename)
	basename = strings.TrimSuffix(basename, ".exe")

	return basename
}
