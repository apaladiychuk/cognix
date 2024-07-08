package utils

import (
	"regexp"
)

// StripFileName removes special characters from the filename and replaces "-" with "_".
func StripFileName(filename string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9.]+`).ReplaceAllString(filename, "")
}
