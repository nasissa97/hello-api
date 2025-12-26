package translation

import (
	"fmt"
	"log"
	"strings"

	"hello-api/handlers/rest"
)

var _ rest.Translator = &RemoteService{}

type RemoteService struct {
	client HelloClient
	cache  map[string]string
}

// HelloClient will call external service.
type HelloClient interface {
	Translate(language, word string) (string, error)
}

func NewRemoteService(client HelloClient) *RemoteService {
	return &RemoteService{
		client: client,
		cache:  make(map[string]string),
	}
}

func (s *RemoteService) Translate(language, word string) string {
	word = strings.ToLower(word)
	language = strings.ToLower(language)

	key := fmt.Sprintf("%s:%s", word, language)

	tr, ok := s.cache[key]
	if ok {
		return tr
	}

	resp, err := s.client.Translate(language, word)
	if err != nil {
		log.Println(err)
		return ""
	}
	s.cache[key] = resp
	return resp
}
