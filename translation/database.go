package translation

import (
	"fmt"

	"hello-api/config"
	"hello-api/handlers/rest"

	"github.com/go-redis/redis"
)

var _ rest.Translator = &Database{}

type Database struct {
	conn *redis.Client
}

func NewDatabaseService(cfg config.Configuration) *Database {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.DatabaseURL, cfg.DatabasePort),
		Password: "",
		DB:       0,
	})
	return &Database{
		conn: rdb,
	}
}

func (s *Database) Close() error {
	return s.conn.Close()
}

func (s *Database) Translate(language string, word string) string {
	out := s.conn.Get(fmt.Sprintf("%s:%s", word, language))
	return out.Val()
}
