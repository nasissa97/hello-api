// Package faas is the entry point for faas.
package faas

import (
	"net/http"

	"hello-api/handlers/rest"
	"hello-api/translation"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	svc := translation.NewStaticService()
	handler := *rest.NewTranslateHandler(svc)

	handler.TranslateHandler(w, r)
}
