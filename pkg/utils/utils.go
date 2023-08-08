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

func GetObjectName(modelType uint8, filename string) string {
	var objectName string
	switch modelType {
	case 0:
		objectName = "checkpoint/" + filename
		break
	case 1:
		objectName = "lora/" + filename
		break
	default:
		objectName = filename
	}
	return objectName
}
