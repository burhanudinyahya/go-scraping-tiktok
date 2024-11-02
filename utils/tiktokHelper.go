package utils

import (
	"regexp"
	"strings"
)

func GetID(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func GetUsername(url string) string {
	parts := strings.Split(url, "/")

	if len(parts) < 4 {
		return ""
	}

	return parts[3]

}

func ExtractTags(text string) []string {
	re := regexp.MustCompile(`#\w+`)
	matches := re.FindAllString(text, -1)

	for i := range matches {
		matches[i] = matches[i][1:] // Skip the first character '#'
	}

	return matches
}
