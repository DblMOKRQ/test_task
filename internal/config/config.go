package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	RestHost string
	RestPort string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &Config{
		Host:     os.Getenv("HOST"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Name:     os.Getenv("NAME"),
		Port:     os.Getenv("PORT"),
		RestHost: os.Getenv("REST_HOST"),
		RestPort: os.Getenv("REST_PORT"),
	}, nil
}
