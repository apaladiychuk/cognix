package utils

import "strings"

// StripFileName removes special characters from the filename and replaces "-" with "_".
func StripFileName(filename string) string {
	return strings.ReplaceAll(strings.ReplaceAll(filename, ":", ""), "-", "_")
}
