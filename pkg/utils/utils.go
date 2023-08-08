package utils

import (
	"fmt"
	"regexp"
)

var RE = regexp.MustCompile(`filename="([^"]+)"`)

func GetDownloadFileName(contentDisposition string) (string, error) {
	match := RE.FindStringSubmatch(contentDisposition)
	if len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("cannot get filename: %s", contentDisposition)
}
