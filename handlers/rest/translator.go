// Package rest handlers for API.
package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"hello-api/translation"
)

const defaultLanguage = "english"

type Resp struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

// TranslateHandler accepts a request then translate hello to language in query.
func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	language := r.URL.Query().Get("language")
	if language == "" {
		language = defaultLanguage
	}
	word := strings.ReplaceAll(r.URL.Path, "/", "")
	translation := translation.Translate(language, word)
	if translation == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp := Resp{
		Language:    language,
		Translation: translation,
	}

	if err := enc.Encode(resp); err != nil {
		panic("unable to encode response")
	}
}
