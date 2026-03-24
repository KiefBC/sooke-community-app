package slug

import (
	"regexp"
	"strings"
)

func GenerateSlug(name string) string {
	normalized := strings.ToLower(name)

	stripChars := []string{"'", "!"}
	replaceChars := []string{"*", "(", ")", "^", "$", "#", "@", "&", "%", "_", " "}

	for _, c := range stripChars {
		normalized = strings.ReplaceAll(normalized, c, "")
	}

	for _, c := range replaceChars {
		normalized = strings.ReplaceAll(normalized, c, "-")
	}

	re := regexp.MustCompile(`-+`)
	normalized = re.ReplaceAllString(normalized, "-")

	normalized = strings.Trim(normalized, "-")

	return normalized
}
