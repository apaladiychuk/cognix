package utils

import "strings"

func StripFileName(filename string) string {
	return strings.ReplaceAll(strings.ReplaceAll(filename, ":", ""), "-", "_")
}
