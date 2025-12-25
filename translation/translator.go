// Package translation is a service to handle translation
package translation

import (
	"strings"
)

// StaticService has data that does not change.
type StaticService struct{}

func NewStaticService() *StaticService {
	return &StaticService{}
}

// Translate converts word from english to language.
func (s *StaticService) Translate(language string, word string) string {
	word = sanitizeInput(word)
	language = sanitizeInput(language)
	if word != "hello" {
		return ""
	}
	switch language {
	case "english":
		return "hello"
	case "finnish":
		return "hei"
	case "german":
		return "hallo"
	case "french":
		return "bonjour"
	default:
		return ""
	}
}

func sanitizeInput(w string) string {
	w = strings.ToLower(w)
	return strings.TrimSpace(w)
}
