package domain

import (
	"regexp"
	"strings"
)

var techAliases = map[string][]string{
	"go":         {"golang", "go", "горутина", "горутины", "goroutine"},
	"python":     {"python", "питон", "django", "fastapi", "pandas"},
	"javascript": {"javascript", "js", "node", "nodejs", "node.js"},
	"react":      {"react", "reactjs", "react.js"},
	"postgres":   {"postgres", "postgresql", "пг"},
	"kubernetes": {"kubernetes", "k8s", "кубер"},
}

func ExtractTags(text string) []string {
	lowerText := strings.ToLower(text)
	var tags []string

	for canonical, aliases := range techAliases {
		for _, alias := range aliases {
			if hasWord(lowerText, alias) {
				tags = append(tags, canonical)
				break
			}
		}
	}
	return tags
}

func CanonicalTags(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	for canonical, aliases := range techAliases {
		for _, a := range aliases {
			if a == input {
				return canonical
			}
		}
	}
	return input
}

func hasWord(text, word string) bool {
	re := regexp.MustCompile(`(^|[^\p{L}])` + regexp.QuoteMeta(word) + `([^\p{L}]|$)`)
	return re.MatchString(text)
}
