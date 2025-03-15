package utils

import (
	"regexp"
	"strings"
)

func CreateSlug(title string) string {
	title = strings.ToLower(title)
	title = strings.ReplaceAll(title, " ", "-")
	reg, _ := regexp.Compile("[^a-z0-9-]+")
	title = reg.ReplaceAllString(title, "")
	return title
}
