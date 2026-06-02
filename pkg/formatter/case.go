package formatter

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchNonAlpha = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = matchNonAlpha.ReplaceAllString(snake, "_")
	return strings.ToLower(strings.Trim(snake, "_"))
}

func ToTitleCase(str string) string {
	text := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	text = matchAllCap.ReplaceAllString(text, "${1}_${2}")
	text = matchNonAlpha.ReplaceAllString(text, " ")

	words := strings.Fields(text)
	for i, word := range words {
		words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
	}

	return strings.Join(words, " ")
}
