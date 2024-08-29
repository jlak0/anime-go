package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	DB      string `json:"db"`
	DB_host string `json:"db_host"`
	DB_user string `json:"db_user"`
	DB_pass string `json:"db_pass"`
	DB_name string `json:"db_name"`
	DB_port int    `json:"db_port"`
}

func LoadConfig(config string) *Config {
	var AppConfig Config
	file, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(bytes, &AppConfig); err != nil {
		log.Fatal(err)
	}
	return &AppConfig
}
