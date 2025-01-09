package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Name string
	Env  string
	Host string
	Port string
	DB   string
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	return &Config{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("ENV"),
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		DB:   os.Getenv("DATABASE_URL"),
	}
}
