// Package config stores objects and functions for configuration.
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Configuration struct {
	Port            string `json:"port"`
	DefaultLanguage string `json:"default_language"`
	LegacyEndpoint  string `json:"legacy_endpoint"`
	DatabaseType    string `json:"database_type"`
	DatabaseURL     string `json:"database_url"`
	DatabasePort    string `json:"database_port"`
}

var defaultConfiguration = Configuration{
	Port:            ":8080",
	DefaultLanguage: "english",
}

func (c *Configuration) LoadFromEnv() {
	if lang := os.Getenv("DEFAULT_LANGUAGE"); lang != "" {
		c.DefaultLanguage = lang
	}
	if port := os.Getenv("PORT"); port != "" {
		c.Port = port
	}
}

func (c *Configuration) ParsePort() {
	if c.Port[0] != ':' {
		c.Port = ":" + c.Port
	}
	if _, err := strconv.Atoi(string(c.Port[1:])); err != nil {
		fmt.Printf("invalid port %s", c.Port)
		c.Port = defaultConfiguration.Port
	}
}

func (c *Configuration) LoadFromJSON(path string) error {
	path = filepath.Clean(path)
	log.Printf("loading configuration from file: %s\n", path)
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("unable to load file %s\n", err.Error())
		return errors.New("unable to load configuration")
	}
	if err := json.Unmarshal(b, c); err != nil {
		log.Printf("unable to parse file: %s\n", err.Error())
		return errors.New("unable to parse configuration")
	}
	if c.Port == "" {
		log.Println("empty port, reverting to default")
		c.Port = defaultConfiguration.DefaultLanguage
	}
	if c.DefaultLanguage == "" {
		log.Println("empty language, reverting to default")
		c.DefaultLanguage = defaultConfiguration.DefaultLanguage
	}

	return nil
}

func LoadConfiguration() Configuration {
	cfgfilePtr := flag.String("config_file", "", "load configuration from a file")
	portPtr := flag.String("port", "", "set port")

	flag.Parse()

	cfg := defaultConfiguration

	if cfgfilePtr != nil && *cfgfilePtr != "" {
		if err := cfg.LoadFromJSON(*cfgfilePtr); err != nil {
			log.Printf("unable to load configuration from json: %s, using default values\n", *cfgfilePtr)
		}
	}

	cfg.LoadFromEnv()

	if portPtr != nil && *portPtr != "" {
		cfg.Port = *portPtr
	}

	cfg.ParsePort()
	return cfg
}
