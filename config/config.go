package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Name string
	Host string
	Port string
}

func MustLoad(cfgPath string) *Config {
	if err := godotenv.Load(cfgPath); err != nil {
		log.Fatal(err)
	}

	return &Config{
		Name: os.Getenv("APP_NAME"),
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}
}
