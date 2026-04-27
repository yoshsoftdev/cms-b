package utils

import (
	"regexp"
	"strings"
)

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = nonAlphaNum.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "untitled"
	}
	return slug
}
